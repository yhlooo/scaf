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

	"github.com/gorilla/websocket"

	metav1 "github.com/yhlooo/scaf/pkg/apis/meta/v1"
	streamv1 "github.com/yhlooo/scaf/pkg/apis/stream/v1"
	"github.com/yhlooo/scaf/pkg/streams"
)

const (
	defaultServer = "http://localhost"
)

// Options 客户端选项
type Options struct {
	// 服务端 URL
	Server string
}

// Complete 将选项补充完整
func (opts *Options) Complete() {
	if opts.Server == "" {
		opts.Server = defaultServer
	}
	opts.Server = strings.TrimSuffix(opts.Server, "/")
}

// Validate 校验选项
func (opts *Options) Validate() error {
	urlObj, err := url.Parse(opts.Server)
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

// New 创建客户端
func New(opts Options) (*Client, error) {
	opts.Complete()
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return &Client{
		opts:       opts,
		httpClient: http.DefaultClient,
		wsDialer:   websocket.DefaultDialer,
	}, nil
}

// Client 客户端
type Client struct {
	opts       Options
	httpClient *http.Client
	wsDialer   *websocket.Dialer
}

// CreateStream 创建流
func (c *Client) CreateStream(ctx context.Context) (*streamv1.Stream, error) {
	ret := &streamv1.Stream{}
	err := c.request(ctx, http.MethodPost, "/v1/streams", nil, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// GetStream 获取流
func (c *Client) GetStream(ctx context.Context, name string) (*streamv1.Stream, error) {
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

// ListStream 列出流
func (c *Client) ListStream(ctx context.Context) (*streamv1.StreamList, error) {
	ret := &streamv1.StreamList{}
	err := c.request(ctx, http.MethodGet, "/v1/streams", nil, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// DeleteStream 删除流
func (c *Client) DeleteStream(ctx context.Context, name string) error {
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
func (c *Client) ConnectStream(ctx context.Context, name string) (streams.Connection, error) {
	conn, resp, connErr := websocket.DefaultDialer.DialContext(ctx, c.opts.Server+"/v1/streams/"+name, nil)
	if connErr == nil {
		return streams.NewWebSocketConnection(conn), nil
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
func (c *Client) request(ctx context.Context, method, uri string, body, resultInto interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		bodyRaw, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body to json error: %w", err)
		}
		bodyReader = bytes.NewReader(bodyRaw)
	}
	uri = "/" + strings.TrimPrefix(uri, "/")
	req, err := http.NewRequestWithContext(ctx, method, c.opts.Server+uri, bodyReader)
	if err != nil {
		return fmt.Errorf("make request error: %w", err)
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
