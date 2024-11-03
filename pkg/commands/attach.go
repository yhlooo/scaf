package commands

import (
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"

	streamv1 "github.com/yhlooo/scaf/pkg/apis/stream/v1"
	clientscommon "github.com/yhlooo/scaf/pkg/clients/common"
	clientsexec "github.com/yhlooo/scaf/pkg/clients/exec"
	"github.com/yhlooo/scaf/pkg/commands/options"
	"github.com/yhlooo/scaf/pkg/streams"
)

// NewAttachCommandWithOptions 基于选项创建 attach 子命令
func NewAttachCommandWithOptions(opts *options.AttachOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attach",
		Short: "Attach to a stream",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := logr.FromContextOrDiscard(ctx)

			// 创建客户端
			client, err := clientsexec.NewTerminal(clientsexec.TerminalOptions{
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
				stream, err := client.Client.CreateStream(ctx, &streamv1.Stream{
					Spec: streamv1.StreamSpec{
						StopPolicy: streamv1.StreamStopPolicy(streams.OnFirstConnectionLeft),
					},
				})
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

			return client.Run(ctx, streamName, os.Stdin, os.Stdout)
		},
	}

	// 绑定选项到命令行
	opts.AddPFlags(cmd.Flags())

	return cmd
}
