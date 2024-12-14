package options

import (
	"fmt"

	"github.com/spf13/pflag"
)

// NewDefaultGlobalOptions 返回默认全局选项
func NewDefaultGlobalOptions() GlobalOptions {
	return GlobalOptions{
		Verbosity: 0,
	}
}

// GlobalOptions 全局选项
type GlobalOptions struct {
	// 日志数量级别（ 0 / 1 / 2 ）
	Verbosity uint32 `json:"verbosity" yaml:"verbosity"`
	// proof 数据导出接口监听地址
	ProfilingAddr string `json:"profilingAddr,omitempty" yaml:"profilingAddr,omitempty"`
	// 将 CPU 性能数据导出到指定文件
	CPUProfile string `json:"cpuProfile,omitempty" yaml:"cpuProfile,omitempty"`
}

// Validate 校验选项是否合法
func (o *GlobalOptions) Validate() error {
	if o.Verbosity > 2 {
		return fmt.Errorf("invalid log verbosity: %d (expected: 0, 1 or 2)", o.Verbosity)
	}
	return nil
}

// AddPFlags 将选项绑定到命令行参数
func (o *GlobalOptions) AddPFlags(fs *pflag.FlagSet) {
	fs.Uint32VarP(&o.Verbosity, "verbose", "v", o.Verbosity, "Number for the log level verbosity (0, 1, or 2)")
	fs.StringVar(&o.ProfilingAddr, "pprof-addr", o.ProfilingAddr, "Address to bind the pprof http server")
	fs.StringVar(&o.CPUProfile, "cpu-profile", o.CPUProfile, "Write a CPU profile to the specified file")
}
