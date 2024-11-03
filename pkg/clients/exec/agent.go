package exec

import (
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/creack/pty"

	"github.com/yhlooo/scaf/pkg/clients/common"
	"github.com/yhlooo/scaf/pkg/streams"
)

// AgentOptions Agent 运行选项
type AgentOptions struct {
	Client common.Options
	TTY    bool
}

// NewAgent 创建 Agent
func NewAgent(opts AgentOptions) (*Agent, error) {
	client, err := common.New(opts.Client)
	if err != nil {
		return nil, err
	}
	return &Agent{
		Client: client,
		tty:    opts.TTY,
	}, nil
}

// Agent exec 代理
type Agent struct {
	*common.Client
	tty bool
}

// Run 与服务端建立连接并运行命令
// 阻塞直到运行结束
func (agent *Agent) Run(ctx context.Context, streamName string, cmd *exec.Cmd) error {
	// 与服务端建立连接
	s := streams.NewBufferedStream()
	if err := s.Start(ctx); err != nil {
		return fmt.Errorf("start local stream error: %w", err)
	}
	defer func() {
		_ = s.Stop(ctx)
	}()
	serverConn, err := agent.Client.ConnectStream(ctx, streamName)
	if err != nil {
		return fmt.Errorf("connect to server error: %w", err)
	}
	if err := s.Join(ctx, serverConn); err != nil {
		return fmt.Errorf("join server connection to local stream error: %w", err)
	}

	if agent.tty {
		// 在 pty 中运行
		ptmx, err := pty.Start(cmd)
		if err != nil {
			return fmt.Errorf("start command error: %w", err)
		}
		if err := s.Join(ctx, streams.NewPassthroughConnection(ptmx, ptmx)); err != nil {
			_ = ptmx.Close()
			return fmt.Errorf("join exec io connection to local stream error: %w", err)
		}
		return cmd.Wait()
	}

	// 设置命令输入输出
	inputR, inputW := io.Pipe()
	outputR, outputW := io.Pipe()
	defer func() {
		_ = inputR.Close()
		_ = outputW.Close()
	}()
	cmd.Stdin = inputR
	cmd.Stdout = outputW
	cmd.Stderr = outputW
	if err := s.Join(ctx, streams.NewPassthroughConnection(outputR, inputW)); err != nil {
		_ = outputR.Close()
		_ = inputW.Close()
		return fmt.Errorf("join exec io connection to local stream error: %w", err)
	}

	// 运行命令
	return cmd.Run()
}
