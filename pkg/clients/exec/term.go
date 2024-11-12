package exec

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-logr/logr"
	"golang.org/x/term"

	"github.com/yhlooo/scaf/pkg/clients/common"
	"github.com/yhlooo/scaf/pkg/streams"
)

const (
	maxReadInputSize = 16 << 10 // 16KiB
)

// NewTerminal 创建 Terminal
func NewTerminal(client *common.Client) *Terminal {
	return &Terminal{c: client}
}

// Terminal exec 终端
type Terminal struct {
	c *common.Client
}

// Client 返回 Terminal 使用的客户端
func (t *Terminal) Client() *common.Client {
	return t.c
}

// WithClient 返回使用指定客户端的 Terminal
func (t *Terminal) WithClient(client *common.Client) *Terminal {
	return &Terminal{
		c: client,
	}
}

// Run 与服务端建立连接并转发输入输出
// 阻塞直到运行结束
func (t *Terminal) Run(ctx context.Context, streamName string, stdin io.Reader, stdout, stderr io.Writer, tty bool) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	logger := logr.FromContextOrDiscard(ctx)

	// 与服务端建立连接
	conn, err := t.c.ConnectStream(ctx, streamName, common.ConnectStreamOptions{
		ConnectionName: "terminal",
	})
	if err != nil {
		return fmt.Errorf("connect to server error: %w", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	// 设置输入流
	var stdinFd *int
	if tty {
		if f, ok := stdin.(*os.File); ok {
			// 将输入流设置为 raw 格式
			stdinFd = new(int)
			*stdinFd = int(f.Fd())
			oldState, err := term.MakeRaw(*stdinFd)
			if err != nil {
				return fmt.Errorf("make input to raw error: %w", err)
			}
			defer func() {
				// 还原输入流
				_ = term.Restore(int(f.Fd()), oldState)
			}()
			// 设置窗口大小
			w, h, err := term.GetSize(*stdinFd)
			if err != nil {
				logger.Error(err, "get terminal size error")
			} else {
				if err := conn.Send(Resize{Height: uint16(h), Width: uint16(w)}.Raw()); err != nil {
					logger.Error(err, "send resize message error")
				}
			}
		}
	}

	// 转发输入输出
	handleConnDone := make(chan struct{})
	go t.handleConn(ctx, handleConnDone, conn, stdout, stderr)
	handleInputDone := make(chan struct{})
	go t.handleInput(ctx, handleInputDone, conn, stdin)

	resizeCh := make(chan os.Signal, 1)
	signal.Notify(resizeCh, syscall.SIGWINCH)
	defer func() {
		signal.Stop(resizeCh)
		close(resizeCh)
	}()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-handleConnDone:
			cancel()
			return nil
		case <-handleInputDone:
			cancel()
			return nil
		case <-resizeCh:
			if stdinFd == nil {
				continue
			}
			w, h, err := term.GetSize(*stdinFd)
			if err != nil {
				logger.Error(err, "get terminal size error")
				continue
			}
			if err := conn.Send(Resize{Height: uint16(h), Width: uint16(w)}.Raw()); err != nil {
				logger.Error(err, "send resize message error")
			}
		}
	}
}

// handleConn 处理连接
func (t *Terminal) handleConn(
	ctx context.Context,
	done chan<- struct{},
	conn streams.Connection,
	stdout, stderr io.Writer,
) {
	defer close(done)
	logger := logr.FromContextOrDiscard(ctx)

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
		case StdoutData:
			if _, err := stdout.Write(m); err != nil {
				logger.Error(err, "write to stdout error")
			}
		case StderrData:
			if _, err := stderr.Write(m); err != nil {
				logger.Error(err, "write to stderr error")
			}
		case ExitCode:
			if m == 0 {
				return
			}
			if _, err := stderr.Write([]byte(fmt.Sprintf("command exit with code: %d\n", m))); err != nil {
				logger.Error(err, "write to stderr error")
			}
			return
		default:
			logger.Info(fmt.Sprintf("unsupported message type: %s, msg: %v", m.Type(), m))
		}
	}
}

// handleInput 处理输入
func (t *Terminal) handleInput(
	ctx context.Context,
	done chan<- struct{},
	conn streams.Connection,
	inputReader io.Reader,
) {
	defer close(done)
	logger := logr.FromContextOrDiscard(ctx)

	tmp := make([]byte, maxReadInputSize)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		n, err := inputReader.Read(tmp)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}

			logger.Error(err, "read input error")
			if err == io.EOF {
				return
			}
			continue
		}

		// 编码消息发送到服务端
		if err := conn.Send(StdinData(tmp[:n]).Raw()); err != nil {
			logger.Error(err, "send message to server error")
		}
	}
}
