package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"

	clientscommon "github.com/yhlooo/scaf/pkg/clients/common"
	clientsexec "github.com/yhlooo/scaf/pkg/clients/exec"
	"github.com/yhlooo/scaf/pkg/commands/options"
)

// NewExecRemoteCommandWithOptions 基于选项创建 exec-remote 子命令
func NewExecRemoteCommandWithOptions(opts *options.ExecRemoteOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec-remote -- COMMAND [ARGS...]",
		Short: "Create remote exec stream and attach to the stream",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := logr.FromContextOrDiscard(ctx)

			// 创建客户端
			client, err := clientscommon.NewClient(opts.Server, opts.Token)
			if err != nil {
				return fmt.Errorf("create client error: %w", err)
			}
			term := clientsexec.NewTerminal(client)

			// 创建流
			stream, err := client.CreateStream(ctx, clientsexec.NewExecStream(args, opts.Input, opts.TTY))
			if err != nil {
				return fmt.Errorf("create stream error: %w", err)
			}
			defer func() {
				if err := client.DeleteStream(ctx, stream.Name); err != nil {
					logger.Error(err, "delete stream error")
				}
			}()
			fmt.Printf("Stream: %s\n", stream.Name)
			execCmd := []string{"scaf", "exec", "-s", opts.Server, "--stream", stream.Name}
			if stream.Status.Token != "" {
				fmt.Printf("Token: %s\n", stream.Status.Token)
				execCmd = append(execCmd, "--token", stream.Status.Token)
				client = client.WithToken(stream.Status.Token)
				term = term.WithClient(client)
			}
			fmt.Printf("Start exec command: %s\n", strings.Join(execCmd, " "))

			return term.Run(ctx, stream, os.Stdin, os.Stdout, os.Stderr)
		},
	}

	// 绑定选项到命令行
	opts.AddPFlags(cmd.Flags())

	return cmd
}
