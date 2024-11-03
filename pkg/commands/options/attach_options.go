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
	// 标准输入是 TTY
	TTY bool `json:"tty,omitempty" yaml:"tty,omitempty"`
}

// AddPFlags 绑定选项到命令行
func (opts *AttachOptions) AddPFlags(fs *pflag.FlagSet) {
	opts.ConnectOptions.AddPFlags(fs)
	fs.BoolVarP(&opts.TTY, "tty", "t", opts.TTY, "Stdin is a TTY")
}
