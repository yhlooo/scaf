package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-logr/logr"
	"github.com/gorilla/websocket"

	authnv1 "github.com/yhlooo/scaf/pkg/apis/authn/v1"
	metav1 "github.com/yhlooo/scaf/pkg/apis/meta/v1"
	streamv1 "github.com/yhlooo/scaf/pkg/apis/stream/v1"
	serverhttp "github.com/yhlooo/scaf/pkg/server/http"
	"github.com/yhlooo/scaf/pkg/streams"
)

const (
	defaultHTTPServer = "http://localhost"
)

// HTTPClientOptions HTTP 客户端选项
type HTTPClientOptions struct {
	// 服务端 URL
	ServerURL string
	// 用于认证的 Token
	Token string
}

// Complete 将选项补充完整
func (opts *HTTPClientOptions) Complete() {
	if opts.ServerURL == "" {
		opts.ServerURL = defaultHTTPServer
	}
	opts.ServerURL = strings.TrimSuffix(opts.ServerURL, "/")
}

// Validate 校验选项
func (opts *HTTPClientOptions) Validate() error {
	urlObj, err := url.Parse(opts.ServerURL)
	if err != nil {
		return fmt.Errorf("invalid server url: %w", err)
	}
	switch urlObj.Scheme {
	case "http", "https":
	default:
		return fmt.Errorf("invalid server url scheme: %q (must be \"https\" or \"http\")", urlObj.Scheme)
	}
	return nil
}

// NewHTTPClient 创建基于 HTTP 的客户端
func NewHTTPClient(opts HTTPClientOptions) (Client, error) {
	opts.Complete()
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return &httpClient{
		opts:       opts,
		httpClient: http.DefaultClient,
		wsDialer:   websocket.DefaultDialer,
	}, nil
}

// httpClient 基于 HTTP 的客户端
type httpClient struct {
	opts       HTTPClientOptions
	httpClient *http.Client
	wsDialer   *websocket.Dialer
}

var _ Client = (*httpClient)(nil)

// Token 返回当前客户端使用的 Token
func (c *httpClient) Token() string {
	return c.opts.Token
}

// WithToken 返回使用指定 Token 的客户端
func (c *httpClient) WithToken(token string) Client {
	opts := c.opts
	opts.Token = token
	return &httpClient{
		opts:       opts,
		httpClient: c.httpClient,
		wsDialer:   c.wsDialer,
	}
}

// Login 登陆获取用户身份返回登陆后的客户端
func (c *httpClient) Login(ctx context.Context, opts LoginOptions) (Client, error) {
	logger := logr.FromContextOrDiscard(ctx)

	if !opts.RenewUser {
		if ret, err := c.CreateSelfSubjectReview(ctx, &authnv1.SelfSubjectReview{}); err == nil {
			// 已经登陆
			logger.V(1).Info(fmt.Sprintf("already login as %q", ret.Status.UserInfo.Username))
			return c, nil
		}
	}
	ret := &authnv1.TokenRequest{}
	err := c.request(ctx, http.MethodPost, "/v1/tokens", &authnv1.TokenRequest{}, ret)
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("login as %q", ret.Name))
	return c.WithToken(ret.Status.Token), nil
}

// CreateSelfSubjectReview 检查自身身份
func (c *httpClient) CreateSelfSubjectReview(
	ctx context.Context,
	review *authnv1.SelfSubjectReview,
) (*authnv1.SelfSubjectReview, error) {
	ret := &authnv1.SelfSubjectReview{}
	err := c.request(ctx, http.MethodPost, "/v1/selfsubjectreviews", review, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// CreateStream 创建流
func (c *httpClient) CreateStream(ctx context.Context, stream *streamv1.Stream) (*streamv1.Stream, error) {
	ret := &streamv1.Stream{}
	err := c.request(ctx, http.MethodPost, "/v1/streams", stream, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// GetStream 获取流
func (c *httpClient) GetStream(ctx context.Context, name string) (*streamv1.Stream, error) {
	if name == "" {
		return nil, fmt.Errorf("stream name must not be empty")
	}
	ret := &streamv1.Stream{}
	err := c.request(ctx, http.MethodGet, "/v1/streams/"+name, nil, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// ListStreams 列出流
func (c *httpClient) ListStreams(ctx context.Context) (*streamv1.StreamList, error) {
	ret := &streamv1.StreamList{}
	err := c.request(ctx, http.MethodGet, "/v1/streams", nil, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// DeleteStream 删除流
func (c *httpClient) DeleteStream(ctx context.Context, name string) error {
	if name == "" {
		return fmt.Errorf("stream name must not be empty")
	}
	ret := &metav1.Status{}
	err := c.request(ctx, http.MethodDelete, "/v1/streams/"+name, nil, ret)
	if err != nil {
		return err
	}
	if ret.Code != http.StatusOK {
		return ret
	}
	return nil
}

// ConnectStream 连接到流
func (c *httpClient) ConnectStream(
	ctx context.Context,
	name string,
	opts ConnectStreamOptions,
) (streams.Connection, error) {
	server := c.opts.ServerURL
	server = strings.Replace(server, "https://", "wss://", 1)
	server = strings.Replace(server, "http://", "ws://", 1)
	header := map[string][]string{
		serverhttp.ConnectionNameHeader: {opts.ConnectionName},
	}
	if c.opts.Token != "" {
		header["Authorization"] = []string{"Bearer " + c.opts.Token}
	}
	conn, resp, connErr := websocket.DefaultDialer.DialContext(ctx, server+"/v1/streams/"+name, header)
	if connErr == nil {
		return streams.NewWebSocketConnection(opts.ConnectionName, conn), nil
	}
	if resp == nil {
		return nil, connErr
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBodyRaw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("connect error: %w, status code: %d", connErr, resp.StatusCode)
	}

	s := &metav1.Status{}
	if err := json.Unmarshal(respBodyRaw, s); err != nil {
		return nil, fmt.Errorf(
			"connect error: %w, status code: %d, body: %s",
			connErr, resp.StatusCode, string(respBodyRaw),
		)
	}

	return nil, s
}

// request 进行一次请求
func (c *httpClient) request(ctx context.Context, method, uri string, body, resultInto interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		bodyRaw, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body to json error: %w", err)
		}
		bodyReader = bytes.NewReader(bodyRaw)
	}
	uri = "/" + strings.TrimPrefix(uri, "/")
	req, err := http.NewRequestWithContext(ctx, method, c.opts.ServerURL+uri, bodyReader)
	if err != nil {
		return fmt.Errorf("make request error: %w", err)
	}
	if c.opts.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.opts.Token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBodyRaw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("read response body error: %w", err)
	}

	if resp.StatusCode >= 400 {
		s := &metav1.Status{}
		if err := json.Unmarshal(respBodyRaw, s); err != nil {
			return fmt.Errorf("unexpected response status code: %d, body: %s", resp.StatusCode, string(respBodyRaw))
		}
		return s
	}

	if err := json.Unmarshal(respBodyRaw, resultInto); err != nil {
		return fmt.Errorf("unmarshal response body from json error: %w", err)
	}

	return nil
}
