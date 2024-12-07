package options

// NewDefaultStreamOptions 创建默认 StreamOptions
func NewDefaultStreamOptions() StreamOptions {
	return StreamOptions{
		Get: StreamGetOptions{
			ClientOptions: NewDefaultClientOptions(),
		},
		List: StreamListOptions{
			ClientOptions: NewDefaultClientOptions(),
		},
		Delete: StreamDeleteOptions{
			ClientOptions: NewDefaultClientOptions(),
		},
	}
}

// StreamOptions stream 子命令选项
type StreamOptions struct {
	// stream get 子命令选项
	Get StreamGetOptions `json:"get,omitempty" yaml:"get,omitempty"`
	// stream list 子命令选项
	List StreamListOptions `json:"list,omitempty" yaml:"list,omitempty"`
	// stream delete 子命令选项
	Delete StreamDeleteOptions `json:"delete,omitempty" yaml:"delete,omitempty"`
}

// StreamGetOptions stream get 子命令选项
type StreamGetOptions struct {
	ClientOptions `yaml:",inline"`
}

// StreamListOptions stream list 子命令选项
type StreamListOptions struct {
	ClientOptions `yaml:",inline"`
}

// StreamDeleteOptions stream delete 子命令选项
type StreamDeleteOptions struct {
	ClientOptions `yaml:",inline"`
}
