package commands

import (
	"encoding/json"
	"fmt"

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
	cmd := &cobra.Command{
		Use: "scaf",
		Short: "Establishes point-to-point streams by pairing and relaying with reverse connections, " +
			"for remote shell, file transfer, etc.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// 校验全局选项
			if err := opts.Global.Validate(); err != nil {
				return err
			}
			// 设置日志
			logger := setLogger(cmd, opts.Global.Verbosity)
			// 输出选项
			optsRaw, _ := json.Marshal(opts)
			logger.V(1).Info(fmt.Sprintf("command: %q, args: %q, options: %s", cmd.Name(), args, string(optsRaw)))
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
