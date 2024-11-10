package exec

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-logr/logr"
	"golang.org/x/term"

	"github.com/yhlooo/scaf/pkg/clients/common"
)

// TerminalOptions Terminal 运行选项
type TerminalOptions struct {
	Client common.Options
	TTY    bool
}

// NewTerminal 创建 Terminal
func NewTerminal(opts TerminalOptions) (*Terminal, error) {
	client, err := common.New(opts.Client)
	if err != nil {
		return nil, err
	}
	return &Terminal{
		Client: client,
		tty:    opts.TTY,
	}, nil
}

// Terminal exec 终端
type Terminal struct {
	*common.Client
	tty bool
}

// WithToken 返回带指定 Token 的客户端
func (t *Terminal) WithToken(token string) *Terminal {
	return &Terminal{
		Client: t.Client.WithToken(token),
		tty:    t.tty,
	}
}

// Run 与服务端建立连接并转发输入输出
// 阻塞直到运行结束
func (t *Terminal) Run(ctx context.Context, streamName string, input io.Reader, output io.Writer) error {
	logger := logr.FromContextOrDiscard(ctx)

	// 与服务端建立连接
	serverConn, err := t.Client.ConnectStream(ctx, streamName, common.ConnectStreamOptions{
		ConnectionName: "terminal",
	})
	if err != nil {
		return fmt.Errorf("connect to server error: %w", err)
	}
	defer func() {
		_ = serverConn.Close()
	}()

	// 设置输入流
	var stdinFd *int
	if t.tty {
		if f, ok := input.(*os.File); ok {
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
				_ = sendResizeANSI(serverConn, h, w)
			}
		}
	}

	// 下行
	downDone := make(chan struct{})
	go func() {
		defer close(downDone)
		if _, err := io.Copy(output, serverConn); err != nil {
			logger.Error(err, "copy from server to terminal error")
		}
	}()
	// 上行
	upDone := make(chan struct{})
	go func() {
		defer close(upDone)
		if _, err := io.Copy(serverConn, input); err != nil {
			logger.Error(err, "copy from terminal to server error")
		}
	}()

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
		case <-upDone:
			return nil
		case <-downDone:
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
			_ = sendResizeANSI(serverConn, h, w)

		}
	}
}

// sendResizeANSI 发送修改窗口大小的 ANSI 序列
// 自定义序列 PM<height>;<width>s
func sendResizeANSI(w io.Writer, height, width int) error {
	_, err := w.Write([]byte(fmt.Sprintf("\x1b^%d;%ds", height, width)))
	return err
}
