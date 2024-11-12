package commands

import (
	"fmt"
	"os/exec"
	"strings"

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
			client, err := clientscommon.New(clientscommon.Options{
				Server: opts.Server,
				Token:  opts.Token,
			})
			if err != nil {
				return fmt.Errorf("create client error: %w", err)
			}
			agent := clientsexec.NewAgent(client)

			// 创建流
			streamName := opts.Stream
			if streamName == "" {
				stream, err := client.CreateStream(ctx, &streamv1.Stream{})
				if err != nil {
					return fmt.Errorf("create stream error: %w", err)
				}
				streamName = stream.Name
				defer func() {
					if err := client.DeleteStream(ctx, streamName); err != nil {
						logger.Error(err, "delete stream error")
					}
				}()
				fmt.Printf("Stream: %s\n", streamName)
				attachCmd := []string{"scaf", "attach", "--stream", streamName}
				if stream.Status.Token != "" {
					fmt.Printf("Token: %s\n", stream.Status.Token)
					attachCmd = append(attachCmd, "--token", stream.Status.Token)
					client = client.WithToken(stream.Status.Token)
					agent = agent.WithClient(client)
				}
				if opts.TTY {
					attachCmd = append(attachCmd, "-t")
				}
				fmt.Printf("Start exec command: %s", strings.Join(attachCmd, " "))
			}

			return agent.Run(ctx, streamName, exec.CommandContext(ctx, args[0], args[1:]...), opts.TTY)
		},
	}

	// 绑定选项到命令行
	opts.AddPFlags(cmd.Flags())

	return cmd
}
