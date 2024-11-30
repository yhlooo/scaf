package options

import "github.com/spf13/pflag"

// NewDefaultReceiveFileOptions 创建默认 ReceiveFileOptions
func NewDefaultReceiveFileOptions() ReceiveFileOptions {
	return ReceiveFileOptions{
		ConnectOptions: NewDefaultConnectOptions(),
	}
}

// ReceiveFileOptions receive-file 子命令选项
type ReceiveFileOptions struct {
	ConnectOptions `yaml:",inline"`
}

// AddPFlags 绑定选项到参数
func (opts *ReceiveFileOptions) AddPFlags(fs *pflag.FlagSet) {
	opts.ConnectOptions.AddPFlags(fs)
}
