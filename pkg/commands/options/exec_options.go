package options

import "github.com/spf13/pflag"

// NewDefaultExecOptions 创建默认 ExecOptions
func NewDefaultExecOptions() ExecOptions {
	return ExecOptions{
		ConnectOptions: NewDefaultConnectOptions(),
		Input:          false,
		TTY:            false,
		Yes:            false,
	}
}

// ExecOptions exec 子命令选项
type ExecOptions struct {
	ConnectOptions `yaml:",inline"`
	// 是否需要开启标准输入流
	Input bool `json:"input,omitempty" yaml:"input,omitempty"`
	// 标准输入是 TTY
	TTY bool `json:"tty,omitempty" yaml:"tty,omitempty"`
	// 是否同意所有二次确认
	Yes bool `json:"yes,omitempty" yaml:"yes,omitempty"`
}

// AddPFlags 绑定选项到命令行
func (opts *ExecOptions) AddPFlags(fs *pflag.FlagSet) {
	opts.ConnectOptions.AddPFlags(fs)
	fs.BoolVarP(&opts.Input, "input", "i", opts.Input, "Enable stdin")
	fs.BoolVarP(&opts.TTY, "tty", "t", opts.TTY, "Stdin is a TTY")
	fs.BoolVarP(&opts.Yes, "yes", "y", opts.Yes, "Skip confirmations and always yes")
}
