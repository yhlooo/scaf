package options

// NewDefaultOptions 创建一个默认运行选项
func NewDefaultOptions() Options {
	return Options{
		Global: NewDefaultGlobalOptions(),

		Serve: NewDefaultServeOptions(),

		Attach:     NewDefaultAttachOptions(),
		Exec:       NewDefaultExecOptions(),
		ExecRemote: NewDefaultExecRemoteOptions(),

		SendFile:    NewDefaultSendFileOptions(),
		ReceiveFile: NewDefaultReceiveFileOptions(),

		Bench: NewDefaultBenchOptions(),

		Stream: NewDefaultStreamOptions(),

		Version: NewDefaultVersionOptions(),
	}
}

// Options pcrctl 运行选项
type Options struct {
	// 全局选项
	Global GlobalOptions `json:"global,omitempty" yaml:"global,omitempty"`

	// serve 子命令选项
	Serve ServeOptions `json:"serve,omitempty" yaml:"serve,omitempty"`

	// attach 子命令选项
	Attach AttachOptions `json:"attach,omitempty" yaml:"attach,omitempty"`
	// exec 子命令选项
	Exec ExecOptions `json:"exec,omitempty" yaml:"exec,omitempty"`
	// exec-remote 子命令选项
	ExecRemote ExecRemoteOptions `json:"execRemote,omitempty" yaml:"execRemote,omitempty"`

	// send-file 子命令选项
	SendFile SendFileOptions `json:"sendFile,omitempty" yaml:"sendFile,omitempty"`
	// receive-file 子命令选项
	ReceiveFile ReceiveFileOptions `json:"receiveFile,omitempty" yaml:"receiveFile,omitempty"`

	// bench 子命令选项
	Bench BenchOptions `json:"bench,omitempty" yaml:"bench,omitempty"`

	// stream 子命令选项
	Stream StreamOptions `json:"stream,omitempty" yaml:"stream,omitempty"`

	// version 子命令选项
	Version VersionOptions `json:"version,omitempty" yaml:"version,omitempty"`
}
