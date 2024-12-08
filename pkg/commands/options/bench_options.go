package options

// NewDefaultBenchOptions 创建默认 BenchOptions
func NewDefaultBenchOptions() BenchOptions {
	return BenchOptions{
		ConnectOptions: NewDefaultConnectOptions(),
	}
}

// BenchOptions bench 子命令选项
type BenchOptions struct {
	ConnectOptions `yaml:",inline"`
}
