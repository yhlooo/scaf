package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	httppprof "net/http/pprof"
	"os"
	"runtime/pprof"

	"github.com/bombsimon/logrusr/v4"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/yhlooo/scaf/pkg/commands/options"
)

// NewScafCommand 创建 scaf 命令
func NewScafCommand() *cobra.Command {
	opts := options.NewDefaultOptions()
	return NewScafCommandWithOptions(&opts)
}

// NewScafCommandWithOptions 创建基于选项的 scaf 命令
func NewScafCommandWithOptions(opts *options.Options) *cobra.Command {
	var cpuprofile *os.File
	cmd := &cobra.Command{
		Use: "scaf",
		Short: "Establishes point-to-point streams by pairing and relaying with reverse connections, " +
			"for remote shell, file transfer, etc.",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// 校验全局选项
			if err := opts.Global.Validate(); err != nil {
				return err
			}
			// 设置日志
			logger := setLogger(cmd, opts.Global.Verbosity)
			// 输出 CPU 性能数据
			if opts.Global.CPUProfile != "" {
				var err error
				cpuprofile, err = os.OpenFile(opts.Global.CPUProfile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
				if err != nil {
					return fmt.Errorf("open cpu profile file %q error: %w", opts.Global.CPUProfile, err)
				}
				if err := pprof.StartCPUProfile(cpuprofile); err != nil {
					return fmt.Errorf("start cpu profile error: %w", err)
				}
			}
			// 启动 pprof http 服务
			if opts.Global.ProfilingAddr != "" {
				if err := startServePprof(cmd.Context(), opts.Global.ProfilingAddr); err != nil {
					return err
				}
			}
			// 输出选项
			optsRaw, _ := json.Marshal(opts)
			logger.V(1).Info(fmt.Sprintf("command: %q, args: %q, options: %s", cmd.Name(), args, string(optsRaw)))
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if cpuprofile != nil {
				pprof.StopCPUProfile()
				if err := cpuprofile.Close(); err != nil {
					return fmt.Errorf("close cpu profile file %q error: %w", opts.Global.CPUProfile, err)
				}
			}
			return nil
		},
	}

	// 将选项绑定到命令行
	opts.Global.AddPFlags(cmd.PersistentFlags())

	// 添加子命令
	cmd.AddCommand(
		NewServeCommandWithOptions(&opts.Serve),

		NewAttachCommandWithOptions(&opts.Attach),
		NewExecCommandWithOptions(&opts.Exec),
		NewExecRemoteCommandWithOptions(&opts.ExecRemote),

		NewSendFileCommandWithOptions(&opts.SendFile),
		NewReceiveFileCommandWithOptions(&opts.ReceiveFile),

		NewBenchCommandWithOptions(&opts.Bench),

		NewStreamCommandWithOptions(&opts.Stream),

		NewVersionCommandWithOptions(&opts.Version),
	)

	return cmd
}

// setLogger 设置命令日志，并返回 logr.Logger
func setLogger(cmd *cobra.Command, verbosity uint32) logr.Logger {
	// 设置日志级别
	logrusLogger := logrus.New()
	switch verbosity {
	case 1:
		logrusLogger.SetLevel(logrus.DebugLevel)
	case 2:
		logrusLogger.SetLevel(logrus.TraceLevel)
	default:
		logrusLogger.SetLevel(logrus.InfoLevel)
	}
	// 将 logger 注入上下文
	logger := logrusr.New(logrusLogger)
	cmd.SetContext(logr.NewContext(cmd.Context(), logger))

	return logger
}

// startServePprof 开始提供 pprof http 服务
func startServePprof(ctx context.Context, addr string) error {
	logger := logr.FromContextOrDiscard(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", httppprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", httppprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", httppprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", httppprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", httppprof.Trace)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen %q for pprof error: %w", addr, err)
	}
	logger.Info(fmt.Sprintf("pprof serve on %q", listener.Addr().String()))

	go func() {
		if err := http.Serve(listener, mux); err != nil {
			logger.Error(err, "serve pprof error")
		}
	}()

	return nil
}
