package options

import "github.com/spf13/pflag"

// NewDefaultServeOptions 创建默认 serve 子命令选项
func NewDefaultServeOptions() ServeOptions {
	return ServeOptions{
		HTTPAddr: ":80",
	}
}

// ServeOptions serve 子命令选项
type ServeOptions struct {
	HTTPAddr string
}

// AddPFlags 绑定选项到参数
func (opts *ServeOptions) AddPFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&opts.HTTPAddr, "listen-http", "l", opts.HTTPAddr, "HTTP listen address")
}
