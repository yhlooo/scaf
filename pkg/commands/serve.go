package commands

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"

	"github.com/yhlooo/scaf/pkg/commands/options"
	"github.com/yhlooo/scaf/pkg/server"
)

// NewServeCommandWithOptions 创建基于选项的 serve 子命令
func NewServeCommandWithOptions(opts *options.ServeOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run scaf server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := logr.FromContextOrDiscard(ctx)

			s := server.NewServer(server.Options{
				HTTPAddr: opts.HTTPAddr,
			})
			if err := s.Start(ctx); err != nil {
				return fmt.Errorf("start server error: %w", err)
			}
			logger.Info(fmt.Sprintf("scaf serve http on %q", s.HTTPAddr().String()))

			// 等待服务结束
			<-s.Done()
			logger.Info("scaf server stopped")

			return nil
		},
	}

	// 绑定选项到命令行参数
	opts.AddPFlags(cmd.Flags())

	return cmd
}
