package options

import (
	"github.com/spf13/pflag"
)

// NewDefaultAttachOptions 创建默认 AttachOptions
func NewDefaultAttachOptions() AttachOptions {
	return AttachOptions{
		ConnectOptions: NewDefaultConnectOptions(),
	}
}

// AttachOptions attach 子命令选项
type AttachOptions struct {
	ConnectOptions `yaml:",inline"`
}

// AddPFlags 绑定选项到命令行
func (opts *AttachOptions) AddPFlags(fs *pflag.FlagSet) {
	opts.ConnectOptions.AddPFlags(fs)
}
