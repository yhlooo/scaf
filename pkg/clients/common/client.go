package common

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-logr/logr"

	authnv1 "github.com/yhlooo/scaf/pkg/apis/authn/v1"
	streamv1 "github.com/yhlooo/scaf/pkg/apis/stream/v1"
	"github.com/yhlooo/scaf/pkg/streams"
)

// Client 客户端
type Client interface {
	// Token 返回当前客户端使用的 Token
	Token() string
	// WithToken 返回使用指定 Token 的客户端
	WithToken(token string) Client
	// Login 登陆获取用户身份返回登陆后的客户端
	Login(ctx context.Context, opts LoginOptions) (Client, error)
	// CreateSelfSubjectReview 检查自身身份
	CreateSelfSubjectReview(ctx context.Context, review *authnv1.SelfSubjectReview) (*authnv1.SelfSubjectReview, error)
	// CreateStream 创建流
	CreateStream(ctx context.Context, stream *streamv1.Stream) (*streamv1.Stream, error)
	// GetStream 获取流
	GetStream(ctx context.Context, name string) (*streamv1.Stream, error)
	// ListStreams 列出流
	ListStreams(ctx context.Context) (*streamv1.StreamList, error)
	// DeleteStream 删除流
	DeleteStream(ctx context.Context, name string) error
	// ConnectStream 连接到流
	ConnectStream(ctx context.Context, name string, opts ConnectStreamOptions) (streams.Connection, error)
}

// LoginOptions 登陆选项
type LoginOptions struct {
	// 是否更换用户
	RenewUser bool
}

// ConnectStreamOptions 连接到流选项
type ConnectStreamOptions struct {
	ConnectionName string
}

// NewClient 创建客户端
func NewClient(server string, token string) (Client, error) {
	urlObj, err := url.Parse(server)
	if err != nil {
		return nil, fmt.Errorf("invalid server url %q: %w", server, err)
	}

	var client Client
	switch urlObj.Scheme {
	case "http", "https":
		client, err = NewHTTPClient(HTTPClientOptions{
			ServerURL: server,
			Token:     token,
		})
	case "grpc":
		client, err = NewGRPCClient(GRPCClientOptions{
			ServerAddress: urlObj.Host,
			Token:         token,
		})
	default:
		return nil, fmt.Errorf("invalid server url %q: unsupported scheme %q", server, urlObj.Scheme)
	}
	if err != nil {
		return nil, err
	}

	tokenFile := ""
	switch runtime.GOOS {
	case "windows":
		tokenFile = os.ExpandEnv("${APPDATA}/scaf/token")
	default:
		tokenFile = os.ExpandEnv("${HOME}/.scaf/token")
	}

	return NewWithPersistentTokenClient(client, tokenFile), nil
}

// NewWithPersistentTokenClient 创建带持久化 Token 的客户端
func NewWithPersistentTokenClient(client Client, tokenFile string) Client {
	return &WithPersistentTokenClient{
		Client:    client,
		tokenFile: tokenFile,
	}
}

// WithPersistentTokenClient 带持久化 Token 的客户端
type WithPersistentTokenClient struct {
	Client
	tokenFile string
}

var _ Client = (*WithPersistentTokenClient)(nil)

// Login 登陆获取用户身份返回登陆后的客户端
func (c *WithPersistentTokenClient) Login(ctx context.Context, opts LoginOptions) (Client, error) {
	logger := logr.FromContextOrDiscard(ctx)

	if opts.RenewUser {
		logger.Info("renew user")
		return c.renewUserLogin(ctx)
	}

	if ret, err := c.Client.CreateSelfSubjectReview(ctx, &authnv1.SelfSubjectReview{}); err == nil {
		// 已经登陆
		logger.Info(fmt.Sprintf("already login as %q", ret.Status.UserInfo.Username))
		return c, nil
	}

	// 读 Token 文件
	token, err := os.ReadFile(c.tokenFile)
	if err != nil {
		if os.IsNotExist(err) {
			return c.renewUserLogin(ctx)
		}
		return nil, fmt.Errorf("read token from %q error: %w", c.tokenFile, err)
	}

	client := c.Client.WithToken(string(token))
	ret, err := client.CreateSelfSubjectReview(ctx, &authnv1.SelfSubjectReview{})
	if err != nil {
		logger.Info(fmt.Sprintf("login with exists token error: %v, renew user", err))
		return c.renewUserLogin(ctx)
	}

	logger.Info(fmt.Sprintf("already login as %q", ret.Status.UserInfo.Username))
	return &WithPersistentTokenClient{Client: client, tokenFile: c.tokenFile}, nil
}

// renewUserLogin 更换用户的登录
func (c *WithPersistentTokenClient) renewUserLogin(ctx context.Context) (Client, error) {
	client, err := c.Client.Login(ctx, LoginOptions{RenewUser: true})
	if err != nil {
		return nil, err
	}

	// 保存 Token
	if err := os.MkdirAll(filepath.Dir(c.tokenFile), 0o700); err != nil {
		return nil, fmt.Errorf("save token to %q error: %w", c.tokenFile, err)
	}
	if err := os.WriteFile(c.tokenFile, []byte(client.Token()), 0o600); err != nil {
		return nil, fmt.Errorf("save token to %q error: %w", c.tokenFile, err)
	}

	return &WithPersistentTokenClient{Client: client, tokenFile: c.tokenFile}, nil
}
