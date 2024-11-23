package options

import "github.com/spf13/pflag"

// NewDefaultExecRemoteOptions 创建默认 ExecRemoteOptions
func NewDefaultExecRemoteOptions() ExecRemoteOptions {
	return ExecRemoteOptions{
		ClientOptions: NewDefaultClientOptions(),
		Input:         false,
		TTY:           false,
	}
}

// ExecRemoteOptions exec-remote 子命令选项
type ExecRemoteOptions struct {
	ClientOptions `yaml:",inline"`
	// 是否需要开启标准输入流
	Input bool `json:"input,omitempty" yaml:"input,omitempty"`
	// 标准输入是 TTY
	TTY bool `json:"tty,omitempty" yaml:"tty,omitempty"`
}

// AddPFlags 绑定选项到命令行
func (opts *ExecRemoteOptions) AddPFlags(fs *pflag.FlagSet) {
	opts.ClientOptions.AddPFlags(fs)
	fs.BoolVarP(&opts.Input, "input", "i", opts.Input, "Enable stdin")
	fs.BoolVarP(&opts.TTY, "tty", "t", opts.TTY, "Stdin is a TTY")
}
