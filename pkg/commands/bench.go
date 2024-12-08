package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"

	streamv1 "github.com/yhlooo/scaf/pkg/apis/stream/v1"
	clientsbench "github.com/yhlooo/scaf/pkg/clients/bench"
	clientscommon "github.com/yhlooo/scaf/pkg/clients/common"
	"github.com/yhlooo/scaf/pkg/commands/options"
)

// NewBenchCommandWithOptions 创建基于选项的 bench 子命令
func NewBenchCommandWithOptions(opts *options.BenchOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bench",
		Aliases: []string{"benchmark"},
		Short:   "Run benchmark",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := logr.FromContextOrDiscard(ctx)

			// 创建客户端
			client, err := clientscommon.NewClient(opts.Server, opts.Token)
			if err != nil {
				return fmt.Errorf("create client error: %w", err)
			}
			if !opts.NoLogin {
				client, err = client.Login(ctx, clientscommon.LoginOptions{RenewUser: opts.RenewUser})
				if err != nil {
					return fmt.Errorf("login error: %w", err)
				}
			}

			if opts.Stream == "" {
				// 创建流并运行测试服务端
				stream, err := client.CreateStream(ctx, &streamv1.Stream{})
				if err != nil {
					return fmt.Errorf("create stream error: %w", err)
				}
				defer func() {
					if err := client.DeleteStream(ctx, stream.Name); err != nil {
						logger.Error(err, "delete stream error")
					}
				}()
				fmt.Printf("Stream: %s\n", stream.Name)
				benchCmd := []string{"scaf", "bench", "-s", opts.Server, "--stream", stream.Name}
				if stream.Status.Token != "" {
					fmt.Printf("Token: %s\n", stream.Status.Token)
					benchCmd = append(benchCmd, "--token", stream.Status.Token)
					client = client.WithToken(stream.Status.Token)
				}
				fmt.Printf("Start benchmark command: %s\n", strings.Join(benchCmd, " "))

				benchServer := clientsbench.NewServer(client)
				return benchServer.Serve(ctx, stream)
			}

			// 获取流
			stream, err := client.GetStream(ctx, opts.Stream)
			if err != nil {
				return fmt.Errorf("get stream %q error: %w", opts.Stream, err)
			}

			// 运行测试客户端
			benchClient := clientsbench.NewClient(client)
			report, err := benchClient.Run(ctx, stream)
			if err != nil {
				return err
			}

			// 展示结果
			raw, err := json.MarshalIndent(report, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal result to json error: %w", err)
			}
			fmt.Println(string(raw))
			return nil
		},
	}

	opts.AddPFlags(cmd.Flags())

	return cmd
}
