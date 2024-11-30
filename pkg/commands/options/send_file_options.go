package options

import (
	"github.com/spf13/pflag"
)

// NewDefaultSendFileOptions 创建默认 SendFileOptions
func NewDefaultSendFileOptions() SendFileOptions {
	return SendFileOptions{
		ClientOptions: NewDefaultClientOptions(),
	}
}

// SendFileOptions send-file 子命令选项
type SendFileOptions struct {
	ClientOptions `yaml:",inline"`
}

// AddPFlags 绑定选项到参数
func (opts *SendFileOptions) AddPFlags(fs *pflag.FlagSet) {
	opts.ClientOptions.AddPFlags(fs)
}
