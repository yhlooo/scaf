package exec

import (
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/creack/pty"
	"github.com/go-logr/logr"

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

// WithToken 返回带指定 Token 的客户端
func (agent *Agent) WithToken(token string) *Agent {
	return &Agent{
		Client: agent.Client.WithToken(token),
		tty:    agent.tty,
	}
}

// Run 与服务端建立连接并运行命令
// 阻塞直到运行结束
func (agent *Agent) Run(ctx context.Context, streamName string, cmd *exec.Cmd) error {
	logger := logr.FromContextOrDiscard(ctx)

	// 与服务端建立连接
	serverConn, err := agent.Client.ConnectStream(ctx, streamName, common.ConnectStreamOptions{
		ConnectionName: "agent",
	})
	if err != nil {
		return fmt.Errorf("connect to server error: %w", err)
	}
	defer func() {
		_ = serverConn.Close()
	}()

	if agent.tty {
		// 在 pty 中运行
		ptmx, err := pty.Start(cmd)
		if err != nil {
			return fmt.Errorf("start command error: %w", err)
		}
		ptyConn := streams.NewPTYConnection("agent", ptmx)
		defer func() {
			_ = ptyConn.Close()
		}()
		ptyConn.InjectLogger(logger)

		// 下行
		go func() {
			if _, err := io.Copy(ptyConn, serverConn); err != nil {
				logger.Error(err, "copy from server to agent error")
			}
			_ = cmd.Process.Kill()
		}()
		// 上行
		go func() {
			if _, err := io.Copy(serverConn, ptyConn); err != nil {
				logger.Error(err, "copy from agent to server error")
			}
			_ = cmd.Process.Kill()
		}()

		return cmd.Wait()
	}

	// 设置命令输入输出
	cmd.Stdin = serverConn
	cmd.Stdout = serverConn
	cmd.Stderr = serverConn

	// 运行命令
	return cmd.Run()
}
