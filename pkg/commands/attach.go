package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	clientscommon "github.com/yhlooo/scaf/pkg/clients/common"
	clientsexec "github.com/yhlooo/scaf/pkg/clients/exec"
	"github.com/yhlooo/scaf/pkg/commands/options"
)

// NewAttachCommandWithOptions 基于选项创建 attach 子命令
func NewAttachCommandWithOptions(opts *options.AttachOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attach",
		Short: "Attach to a stream",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// 创建客户端
			client, err := clientscommon.NewClient(opts.Server, opts.Token)
			if err != nil {
				return fmt.Errorf("create client error: %w", err)
			}
			term := clientsexec.NewTerminal(client)

			// 获取流
			stream, err := client.GetStream(ctx, opts.Stream)
			if err != nil {
				return fmt.Errorf("get stream %q error: %w", opts.Stream, err)
			}

			return term.Run(ctx, stream, os.Stdin, os.Stdout, os.Stderr)
		},
	}

	// 绑定选项到命令行
	opts.AddPFlags(cmd.Flags())

	return cmd
}
