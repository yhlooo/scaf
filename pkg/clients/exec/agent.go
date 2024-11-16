package exec

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/go-logr/logr"

	"github.com/yhlooo/scaf/pkg/clients/common"
	"github.com/yhlooo/scaf/pkg/streams"
)

const (
	maxReadOutputSize = 16 << 10 // 16KiB
)

// NewAgent 创建 Agent
func NewAgent(client common.Client) *Agent {
	return &Agent{c: client}
}

// Agent exec 代理
type Agent struct {
	c common.Client
}

// Client 返回 Agent 使用的客户端
func (agent *Agent) Client() common.Client {
	return agent.c
}

// WithClient 返回使用指定客户端的 Agent
func (agent *Agent) WithClient(client common.Client) *Agent {
	return &Agent{
		c: client,
	}
}

// Run 与服务端建立连接并运行命令
// 阻塞直到运行结束
func (agent *Agent) Run(ctx context.Context, streamName string, cmd *exec.Cmd, tty bool) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	logger := logr.FromContextOrDiscard(ctx)

	// 与服务端建立连接
	conn, err := agent.c.ConnectStream(ctx, streamName, common.ConnectStreamOptions{
		ConnectionName: "agent",
	})
	if err != nil {
		return fmt.Errorf("connect to server error: %w", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	var inputWriter io.Writer
	var outputReader io.Reader
	var errorReader io.Reader

	// 启动命令
	if tty {
		// 在 pty 中运行
		ptmx, err := pty.Start(cmd)
		if err != nil {
			return fmt.Errorf("start command error: %w", err)
		}
		defer func() {
			if err := ptmx.Close(); err != nil {
				logger.Error(err, "close pty error")
			}
		}()
		inputWriter = ptmx
		outputReader = ptmx
	} else {
		// 设置命令输入输出管道
		stdinR, stdinW := io.Pipe()
		stdoutR, stdoutW := io.Pipe()
		stderrR, stderrW := io.Pipe()
		defer func() {
			_ = stdinR.Close()
			_ = stdinW.Close()
			_ = stdoutR.Close()
			_ = stdoutW.Close()
			_ = stderrR.Close()
			_ = stderrW.Close()
		}()
		inputWriter = stdinW
		outputReader = stdoutR
		errorReader = stderrR

		cmd.Stdin = stdinR
		cmd.Stdout = stdoutW
		cmd.Stderr = stderrW
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("start command error: %w", err)
		}
	}

	// 转发输入输出
	handleConnDone := make(chan struct{})
	go agent.handleConn(ctx, handleConnDone, conn, inputWriter)
	handleStdoutDone := make(chan struct{})
	go agent.handleOutput(ctx, handleStdoutDone, conn, outputReader, false)
	if errorReader != nil {
		handleStderrDone := make(chan struct{})
		go agent.handleOutput(ctx, handleStderrDone, conn, errorReader, true)
	}
	go func() {
		select {
		case <-handleConnDone:
			cancel()
		case <-handleStdoutDone:
			cancel()
		case <-ctx.Done():
		}
		_ = cmd.Process.Kill()
	}()

	// 等待命令结束
	err = cmd.Wait()
	var sendErr error
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			sendErr = conn.Send(ExitCode(exitErr.ExitCode()).Raw())
		} else {
			sendErr = conn.Send(ExitCode(255).Raw())
		}
	} else {
		sendErr = conn.Send(ExitCode(0).Raw())
	}
	if sendErr != nil {
		logger.Error(sendErr, "send exit code error")
	}

	return err
}

// handleConn 处理连接
func (agent *Agent) handleConn(
	ctx context.Context,
	done chan<- struct{},
	conn streams.Connection,
	stdinWriter io.Writer,
) {
	defer close(done)
	logger := logr.FromContextOrDiscard(ctx)

	stdinFile, _ := stdinWriter.(*os.File)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		data, err := conn.Receive()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			logger.Error(err, "receive from server error")
			if errors.Is(err, streams.ErrConnectionClosed) {
				return
			}
			continue
		}

		// 解析消息
		msg, err := ParseMessage(data)
		if err != nil {
			logger.Error(err, "parse message error")
			continue
		}

		// 分类处理
		switch m := msg.(type) {
		case StdinData:
			if _, err := stdinWriter.Write(m); err != nil {
				logger.Error(err, "write to stdin error")
			}
		case Resize:
			if stdinFile == nil {
				logger.Info(fmt.Sprintf("not support resize pty, msg: %v", m))
			}
			if err := pty.Setsize(stdinFile, &pty.Winsize{Rows: m.Height, Cols: m.Width}); err != nil {
				logger.Error(err, "set pty size error")
			}
		default:
			logger.Info(fmt.Sprintf("unsupported message type: %s, msg: %v", m.Type(), m))
		}
	}
}

// handleOutput 处理输出
func (agent *Agent) handleOutput(
	ctx context.Context,
	done chan<- struct{},
	conn streams.Connection,
	outputReader io.Reader,
	stderr bool,
) {
	defer close(done)
	logger := logr.FromContextOrDiscard(ctx)

	tmp := make([]byte, maxReadOutputSize)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		n, err := outputReader.Read(tmp)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}

			logger.Error(err, "read output error")
			if err == io.EOF {
				return
			}
			continue
		}

		// 编码消息发送到服务端
		var msg Message
		if stderr {
			msg = StderrData(tmp[:n])
		} else {
			msg = StdoutData(tmp[:n])
		}
		if err := conn.Send(msg.Raw()); err != nil {
			logger.Error(err, "send message to server error")
		}
	}
}
