package commands

import (
	"fmt"
	"os/exec"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"

	streamv1 "github.com/yhlooo/scaf/pkg/apis/stream/v1"
	clientscommon "github.com/yhlooo/scaf/pkg/clients/common"
	clientsexec "github.com/yhlooo/scaf/pkg/clients/exec"
	"github.com/yhlooo/scaf/pkg/commands/options"
)

// NewExecCommandWithOptions 基于选项创建 exec 子命令
func NewExecCommandWithOptions(opts *options.ExecOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec -- COMMAND [ARGS...]",
		Short: "Execute command and forward input and output through stream",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("the command to be executed must be specified")
			}

			ctx := cmd.Context()
			logger := logr.FromContextOrDiscard(ctx)

			// 创建客户端
			client, err := clientsexec.NewAgent(clientsexec.AgentOptions{
				Client: clientscommon.Options{
					Server: opts.Server,
				},
				TTY: opts.TTY,
			})
			if err != nil {
				return fmt.Errorf("create client error: %w", err)
			}

			// 创建流
			streamName := opts.Stream
			if streamName == "" {
				stream, err := client.Client.CreateStream(ctx, &streamv1.Stream{})
				if err != nil {
					return fmt.Errorf("create stream error: %w", err)
				}
				streamName = stream.Name
				fmt.Printf("Stream: %s\n", streamName)
				defer func() {
					if err := client.DeleteStream(ctx, streamName); err != nil {
						logger.Error(err, "delete stream error")
					}
				}()
			}

			return client.Run(ctx, streamName, exec.CommandContext(ctx, args[0], args[1:]...))
		},
	}

	// 绑定选项到命令行
	opts.AddPFlags(cmd.Flags())

	return cmd
}
