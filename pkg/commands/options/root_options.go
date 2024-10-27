package options

// NewDefaultOptions 创建一个默认运行选项
func NewDefaultOptions() Options {
	return Options{
		Global:  NewDefaultGlobalOptions(),
		Serve:   NewDefaultServeOptions(),
		Version: NewDefaultVersionOptions(),
	}
}

// Options pcrctl 运行选项
type Options struct {
	// 全局选项
	Global GlobalOptions `json:"global,omitempty" yaml:"global,omitempty"`

	// serve 子命令选项
	Serve ServeOptions `json:"serve,omitempty" yaml:"serve,omitempty"`
	// version 子命令选项
	Version VersionOptions `json:"version,omitempty" yaml:"version,omitempty"`
}
