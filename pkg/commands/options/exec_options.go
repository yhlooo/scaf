package options

import "github.com/spf13/pflag"

// NewDefaultExecOptions 创建默认 ExecOptions
func NewDefaultExecOptions() ExecOptions {
	return ExecOptions{
		ConnectOptions: NewDefaultConnectOptions(),
		TTY:            false,
	}
}

// ExecOptions exec 子命令选项
type ExecOptions struct {
	ConnectOptions `yaml:",inline"`
	// 标准输入是 TTY
	TTY bool `json:"tty,omitempty" yaml:"tty,omitempty"`
}

// AddPFlags 绑定选项到命令行
func (opts *ExecOptions) AddPFlags(fs *pflag.FlagSet) {
	opts.ConnectOptions.AddPFlags(fs)
	fs.BoolVarP(&opts.TTY, "tty", "t", opts.TTY, "Stdin is a TTY")
}
