package commands

import (
	"encoding/json"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"

	clientscommon "github.com/yhlooo/scaf/pkg/clients/common"
	"github.com/yhlooo/scaf/pkg/commands/options"
)

// NewStreamCommandWithOptions 创建基于选项的 stream 子命令
func NewStreamCommandWithOptions(opts *options.StreamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stream",
		Short: "Manage streams",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(
		NewStreamGetCommandWithOptions(&opts.Get),
		NewStreamListCommandWithOptions(&opts.List),
		NewStreamDeleteCommandWithOptions(&opts.Delete),
	)
	return cmd
}

// NewStreamGetCommandWithOptions 创建基于选项的 stream get 子命令
func NewStreamGetCommandWithOptions(opts *options.StreamGetOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get STREAM_NAME",
		Short: "Get stream info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

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

			stream, err := client.GetStream(ctx, args[0])
			if err != nil {
				return err
			}

			raw, err := json.MarshalIndent(stream, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal result to json error: %w", err)
			}

			fmt.Println(string(raw))
			return nil
		},
	}

	// 绑定选项到命令行参数
	opts.AddPFlags(cmd.Flags())

	return cmd
}

// NewStreamListCommandWithOptions 创建基于选项的 stream list 子命令
func NewStreamListCommandWithOptions(opts *options.StreamListOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List streams",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

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

			streamList, err := client.ListStreams(ctx)
			if err != nil {
				return err
			}

			for _, stream := range streamList.Items {
				fmt.Println(stream.Name)
			}

			return nil
		},
	}

	// 绑定选项到命令行参数
	opts.AddPFlags(cmd.Flags())

	return cmd
}

// NewStreamDeleteCommandWithOptions 创建基于选项的 stream delete子命令
func NewStreamDeleteCommandWithOptions(opts *options.StreamDeleteOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete STREAM_NAME",
		Short: "Delete stream",
		Args:  cobra.ExactArgs(1),
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

			if err := client.DeleteStream(ctx, args[0]); err != nil {
				return err
			}

			logger.Info(fmt.Sprintf("stream %q deleted", args[0]))
			return nil
		},
	}

	// 绑定选项到命令行参数
	opts.AddPFlags(cmd.Flags())

	return cmd
}
