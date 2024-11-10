package options

import "github.com/spf13/pflag"

// NewDefaultConnectOptions 创建默认 ConnectOptions
func NewDefaultConnectOptions() ConnectOptions {
	return ConnectOptions{
		Server: "http://localhost",
		Stream: "",
		Token:  "",
	}
}

// ConnectOptions 建立连接公共选项
type ConnectOptions struct {
	// 服务端地址
	Server string `json:"server,omitempty" yaml:"server,omitempty"`
	// 连接的流名
	Stream string `json:"stream,omitempty" yaml:"stream,omitempty"`
	// 用于认证的 Token
	Token string `json:"token,omitempty" yaml:"token,omitempty"`
}

// AddPFlags 绑定选项到命令行
func (opts *ConnectOptions) AddPFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&opts.Server, "server", "s", opts.Server, "Server address")
	fs.StringVar(&opts.Stream, "stream", opts.Stream, "Stream name connect to")
	fs.StringVar(&opts.Token, "token", opts.Token, "Token")
}
