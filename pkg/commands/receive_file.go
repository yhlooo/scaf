package commands

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"

	clientscp "github.com/yhlooo/scaf/pkg/clients/cp"
	"github.com/yhlooo/scaf/pkg/commands/options"
)

// NewReceiveFileCommandWithOptions 基于选项创建 receive-file 子命令
func NewReceiveFileCommandWithOptions(opts *options.ReceiveFileOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "receive-file [PATH]",
		Short: "Receive file from stream",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := logr.FromContextOrDiscard(ctx)

			// 创建客户端
			client, err := opts.NewClient(ctx)
			if err != nil {
				return fmt.Errorf("create client error: %w", err)
			}
			cpClient := clientscp.New(client)

			// 获取流
			stream, err := client.GetStream(ctx, opts.Stream)
			if err != nil {
				return fmt.Errorf("get stream %q error: %w", opts.Stream, err)
			}

			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			target, err := cpClient.Receive(ctx, stream, path)
			if err != nil {
				return err
			}
			logger.Info(fmt.Sprintf("received %q", target))
			return nil
		},
	}

	// 绑定选项到命令行参数
	opts.AddPFlags(cmd.Flags())

	return cmd
}
