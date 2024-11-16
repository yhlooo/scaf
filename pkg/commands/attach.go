package commands

import (
	"fmt"
	"os"
	"strings"

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
			var client clientscommon.Client
			var err error
			switch {
			case opts.HTTP:
				client, err = clientscommon.NewHTTPClient(clientscommon.HTTPClientOptions{
					ServerURL: opts.Server,
					Token:     opts.Token,
				})
			case opts.GRPC:
				client, err = clientscommon.NewGRPCClient(clientscommon.GRPCClientOptions{
					ServerAddress: opts.Server,
					Token:         opts.Token,
				})
			default:
				// TODO: 应该自动判断服务端支持什么模式，这里暂时使用 gRPC
				client, err = clientscommon.NewGRPCClient(clientscommon.GRPCClientOptions{
					ServerAddress: opts.Server,
					Token:         opts.Token,
				})
			}
			if err != nil {
				return fmt.Errorf("create client error: %w", err)
			}
			term := clientsexec.NewTerminal(client)

			// 创建流
			streamName := opts.Stream
			if streamName == "" {
				stream, err := client.CreateStream(ctx, &streamv1.Stream{
					Spec: streamv1.StreamSpec{
						StopPolicy: streamv1.StreamStopPolicy(streams.OnFirstConnectionLeft),
					},
				})
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
				execCmd := []string{"scaf", "exec", "--stream", streamName}
				if stream.Status.Token != "" {
					fmt.Printf("Token: %s\n", stream.Status.Token)
					execCmd = append(execCmd, "--token", stream.Status.Token)
					client = client.WithToken(stream.Status.Token)
					term = term.WithClient(client)
				}
				if opts.TTY {
					execCmd = append(execCmd, "-t")
				}
				execCmd = append(execCmd, "--", "COMMAND", "ARG...")
				fmt.Printf("Start exec command: %s", strings.Join(execCmd, " "))
			}

			return term.Run(ctx, streamName, os.Stdin, os.Stdout, os.Stderr, opts.TTY)
		},
	}

	// 绑定选项到命令行
	opts.AddPFlags(cmd.Flags())

	return cmd
}
