package commands

import (
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"

	streamv1 "github.com/yhlooo/scaf/pkg/apis/stream/v1"
	clientscommon "github.com/yhlooo/scaf/pkg/clients/common"
	clientscp "github.com/yhlooo/scaf/pkg/clients/cp"
	"github.com/yhlooo/scaf/pkg/commands/options"
)

// NewSendFileCommandWithOptions 基于选项创建 send-file 子命令
func NewSendFileCommandWithOptions(opts *options.SendFileOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-file PATH",
		Short: "Send file to stream",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := logr.FromContextOrDiscard(ctx)

			// 创建客户端
			client, err := clientscommon.NewClient(opts.Server, opts.Token)
			if err != nil {
				return fmt.Errorf("create client error: %w", err)
			}
			cpClient := clientscp.New(client)

			// 创建流
			stream, err := client.CreateStream(ctx, &streamv1.Stream{
				Spec: streamv1.StreamSpec{StopPolicy: streamv1.OnFirstConnectionLeft},
			})
			if err != nil {
				return fmt.Errorf("create stream error: %w", err)
			}
			defer func() {
				if err := client.DeleteStream(ctx, stream.Name); err != nil {
					logger.Error(err, "delete stream error")
				}
			}()
			fmt.Printf("Stream: %s\n", stream.Name)
			recvCmd := []string{"scaf", "receive-file", "-s", opts.Server, "--stream", stream.Name}
			if stream.Status.Token != "" {
				fmt.Printf("Token: %s\n", stream.Status.Token)
				recvCmd = append(recvCmd, "--token", stream.Status.Token)
				client = client.WithToken(stream.Status.Token)
				cpClient = cpClient.WithClient(client)
			}
			fmt.Printf("Receive file command: %s\n", strings.Join(recvCmd, " "))

			if err := cpClient.Send(ctx, stream, args[0]); err != nil {
				return err
			}
			logger.Info("done")
			return nil
		},
	}

	// 绑定选项到命令行参数
	opts.AddPFlags(cmd.Flags())

	return cmd
}
