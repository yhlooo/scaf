package options

import "github.com/spf13/pflag"

// NewDefaultConnectOptions 创建默认 ConnectOptions
func NewDefaultConnectOptions() ConnectOptions {
	return ConnectOptions{
		Server: "http://localhost",
		Stream: "",
	}
}

// ConnectOptions 建立连接公共选项
type ConnectOptions struct {
	// 服务端地址
	Server string `json:"server,omitempty" yaml:"server,omitempty"`
	// 连接的流地址
	Stream string `json:"stream,omitempty" yaml:"stream,omitempty"`
}

// AddPFlags 绑定选项到命令行
func (opts *ConnectOptions) AddPFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&opts.Server, "server", "s", opts.Server, "Server address")
	fs.StringVar(&opts.Stream, "stream", opts.Stream, "Stream name connect to")
}
