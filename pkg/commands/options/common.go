package options

import (
	"context"
	"fmt"

	"github.com/spf13/pflag"

	clientscommon "github.com/yhlooo/scaf/pkg/clients/common"
)

// NewDefaultClientOptions 创建默认 ClientOptions
func NewDefaultClientOptions() ClientOptions {
	return ClientOptions{
		Server:    "grpc://localhost:9443",
		Token:     "",
		NoLogin:   false,
		RenewUser: false,
	}
}

// ClientOptions 客户端选项
type ClientOptions struct {
	// 服务端地址
	Server string `json:"server,omitempty" yaml:"server,omitempty"`
	// 用于认证的 Token
	Token string `json:"token,omitempty" yaml:"token,omitempty"`
	// 不登陆，使用匿名用户访问
	NoLogin bool `json:"noLogin,omitempty" yaml:"noLogin,omitempty"`
	// 始终使用新用户登录
	RenewUser bool `json:"renewUser,omitempty" yaml:"renewUser,omitempty"`
	// 是否对传输数据进行压缩
	Compress bool `json:"compress,omitempty" yaml:"compress,omitempty"`
}

// AddPFlags 绑定选项到命令行
func (opts *ClientOptions) AddPFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&opts.Server, "server", "s", opts.Server, "Server address")
	fs.StringVar(&opts.Token, "token", opts.Token, "Token")
	fs.BoolVar(&opts.NoLogin, "no-login", opts.NoLogin, "Do not login and access anonymously")
	fs.BoolVar(&opts.RenewUser, "renew-user", opts.RenewUser, "Renew user")
	fs.BoolVar(&opts.Compress, "compress", opts.Compress, "Compress the transport stream")
}

// NewClient 基于选项创建客户端
func (opts *ClientOptions) NewClient(ctx context.Context) (clientscommon.Client, error) {
	client, err := clientscommon.NewClient(clientscommon.ClientOptions{
		Server:   opts.Server,
		Token:    opts.Token,
		Compress: opts.Compress,
	})
	if err != nil {
		return nil, err
	}
	if !opts.NoLogin {
		client, err = client.Login(ctx, clientscommon.LoginOptions{RenewUser: opts.RenewUser})
		if err != nil {
			return nil, fmt.Errorf("login error: %w", err)
		}
	}
	return client, nil
}

// NewDefaultConnectOptions 创建默认 ConnectOptions
func NewDefaultConnectOptions() ConnectOptions {
	return ConnectOptions{
		ClientOptions: NewDefaultClientOptions(),
		Stream:        "",
	}
}

// ConnectOptions 建立连接公共选项
type ConnectOptions struct {
	ClientOptions `yaml:",inline"`
	// 连接的流名
	Stream string `json:"stream,omitempty" yaml:"stream,omitempty"`
}

// AddPFlags 绑定选项到命令行
func (opts *ConnectOptions) AddPFlags(fs *pflag.FlagSet) {
	opts.ClientOptions.AddPFlags(fs)
	fs.StringVar(&opts.Stream, "stream", opts.Stream, "Stream name connect to")
}
