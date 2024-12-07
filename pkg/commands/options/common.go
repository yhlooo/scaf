package options

import "github.com/spf13/pflag"

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
}

// AddPFlags 绑定选项到命令行
func (opts *ClientOptions) AddPFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&opts.Server, "server", "s", opts.Server, "Server address")
	fs.StringVar(&opts.Token, "token", opts.Token, "Token")
	fs.BoolVar(&opts.NoLogin, "no-login", opts.NoLogin, "Do not login and access anonymously")
	fs.BoolVar(&opts.RenewUser, "renew-user", opts.RenewUser, "Renew user")
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
