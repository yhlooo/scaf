package commands

import (
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"

	streamv1 "github.com/yhlooo/scaf/pkg/apis/stream/v1"
	clientsexec "github.com/yhlooo/scaf/pkg/clients/exec"
	"github.com/yhlooo/scaf/pkg/commands/options"
)

// NewExecCommandWithOptions 基于选项创建 exec 子命令
func NewExecCommandWithOptions(opts *options.ExecOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec [-- COMMAND [ARGS...]]",
		Short: "Execute command and forward input and output through stream",
		Example: `# Create a new stream
scaf exec [-i] [-t] -s SERVER -- COMMAND [ARGS...]

# Join an existing stream
scaf exec -s SERVER --stream STREAM --token TOKEN`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := logr.FromContextOrDiscard(ctx)

			// 创建客户端
			client, err := opts.NewClient(ctx)
			if err != nil {
				return fmt.Errorf("create client error: %w", err)
			}
			agent := clientsexec.NewAgent(client)

			var stream *streamv1.Stream
			if opts.Stream != "" {
				// 获取流
				stream, err = client.GetStream(ctx, opts.Stream)
				if err != nil {
					return fmt.Errorf("get stream %q error: %w", opts.Stream, err)
				}
				// 解析执行信息
				command, input, tty, err := clientsexec.GetExecOptions(stream)
				if err != nil {
					return fmt.Errorf("get exec options error: %w", err)
				}
				fmt.Printf("Command: %q\n", command)
				fmt.Printf("Input:   %t\n", input)
				fmt.Printf("TTY:     %t\n", tty)
				// 二次确认
				if !opts.Yes {
					fmt.Print("Continue? (Y/n): ")
					confirm := ""
					_, _ = fmt.Scanln(&confirm)
					if confirm != "Y" {
						fmt.Println("abort")
						return nil
					}
				}
			} else {
				// 创建流
				stream = clientsexec.NewExecStream(args, opts.Input, opts.TTY)
				newStream, err := client.CreateStream(ctx, stream)
				if err != nil {
					return fmt.Errorf("create stream error: %w", err)
				}
				defer func() {
					if err := client.DeleteStream(ctx, stream.Name); err != nil {
						logger.Error(err, "delete stream error")
					}
				}()
				stream.Name = newStream.Name // 不直接使用 newStream 是避免被恶意的服务端返回内容篡改实际执行的命令
				stream.UID = newStream.UID

				fmt.Printf("Stream: %s\n", newStream.Name)
				attachCmd := []string{"scaf", "attach", "-s", opts.Server, "--stream", newStream.Name}
				if newStream.Status.Token != "" {
					fmt.Printf("Token: %s\n", newStream.Status.Token)
					attachCmd = append(attachCmd, "--token", newStream.Status.Token)
					client = client.WithToken(newStream.Status.Token)
					agent = agent.WithClient(client)
				}
				fmt.Printf("Start exec command: %s\n", strings.Join(attachCmd, " "))
			}

			return agent.Run(ctx, stream)
		},
	}

	// 绑定选项到命令行
	opts.AddPFlags(cmd.Flags())

	return cmd
}
