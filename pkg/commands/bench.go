package commands

import (
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"

	streamv1 "github.com/yhlooo/scaf/pkg/apis/stream/v1"
	clientsbench "github.com/yhlooo/scaf/pkg/clients/bench"
	"github.com/yhlooo/scaf/pkg/commands/options"
	"github.com/yhlooo/scaf/pkg/utils/units"
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
			client, err := opts.NewClient(ctx)
			if err != nil {
				return fmt.Errorf("create client error: %w", err)
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
			fmt.Println()
			fmt.Println("Result:")
			fmt.Println()
			fmt.Println("Ping:")
			fmt.Printf(
				"  Time:   %s\t(lost: %.2f%%)\n",
				report.Ping.RoundTripTime, report.Ping.LossRate*100,
			)
			fmt.Println("Read Only:")
			fmt.Printf(
				"  Read:   %sB/s\t(read: %sB, packages: %s, lost: %.2f%%)\n",
				units.NewIECValue(int64(report.ReadOnly.Throughput)).RoundString(2),
				units.NewIECValue(int64(report.ReadOnly.Size)).RoundString(2),
				units.NewSIValue(int64(report.ReadOnly.Packages)).RoundString(2),
				report.ReadOnly.LossRate*100,
			)
			fmt.Println("Write Only:")
			fmt.Printf(
				"  Write:  %sB/s\t(write: %sB, packages: %s, lost: %.2f%%)\n",
				units.NewIECValue(int64(report.WriteOnly.Throughput)).RoundString(2),
				units.NewIECValue(int64(report.WriteOnly.Size)).RoundString(2),
				units.NewSIValue(int64(report.WriteOnly.Packages)).RoundString(2),
				report.WriteOnly.LossRate*100,
			)
			fmt.Println("Read and Write:")
			fmt.Printf(
				"  Read:   %sB/s\t(read: %sB, packages: %s, lost: %.2f%%)\n",
				units.NewIECValue(int64(report.ReadWrite.Read.Throughput)).RoundString(2),
				units.NewIECValue(int64(report.ReadWrite.Read.Size)).RoundString(2),
				units.NewSIValue(int64(report.ReadWrite.Read.Packages)).RoundString(2),
				report.ReadWrite.Read.LossRate*100,
			)
			fmt.Printf(
				"  Write:  %sB/s\t(write: %sB, packages: %s, lost: %.2f%%)\n",
				units.NewIECValue(int64(report.ReadWrite.Write.Throughput)).RoundString(2),
				units.NewIECValue(int64(report.ReadWrite.Write.Size)).RoundString(2),
				units.NewSIValue(int64(report.ReadWrite.Write.Packages)).RoundString(2),
				report.ReadWrite.Write.LossRate*100,
			)
			return nil
		},
	}

	opts.AddPFlags(cmd.Flags())

	return cmd
}
