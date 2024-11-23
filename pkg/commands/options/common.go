package options

import "github.com/spf13/pflag"

// NewDefaultClientOptions 创建默认 ClientOptions
func NewDefaultClientOptions() ClientOptions {
	return ClientOptions{
		Server: "grpc://localhost:9443",
		Token:  "",
	}
}

// ClientOptions 客户端选项
type ClientOptions struct {
	// 服务端地址
	Server string `json:"server,omitempty" yaml:"server,omitempty"`
	// 用于认证的 Token
	Token string `json:"token,omitempty" yaml:"token,omitempty"`
}

// AddPFlags 绑定选项到命令行
func (opts *ClientOptions) AddPFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&opts.Server, "server", "s", opts.Server, "Server address")
	fs.StringVar(&opts.Token, "token", opts.Token, "Token")
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
