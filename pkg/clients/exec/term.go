package exec

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/go-logr/logr"
	"golang.org/x/term"

	"github.com/yhlooo/scaf/pkg/clients/common"
	"github.com/yhlooo/scaf/pkg/streams"
)

// TerminalOptions Terminal 运行选项
type TerminalOptions struct {
	Client common.Options
}

// NewTerminal 创建 Terminal
func NewTerminal(opts TerminalOptions) (*Terminal, error) {
	client, err := common.New(opts.Client)
	if err != nil {
		return nil, err
	}
	return &Terminal{
		Client: client,
	}, nil
}

// Terminal exec 终端
type Terminal struct {
	*common.Client
}

// Run 与服务端建立连接并转发输入输出
// 阻塞直到运行结束
func (t *Terminal) Run(ctx context.Context, streamName string, input io.Reader, output io.Writer) error {
	logger := logr.FromContextOrDiscard(ctx)

	// 与服务端建立连接
	s := streams.NewBufferedStream()
	if err := s.Start(ctx); err != nil {
		return fmt.Errorf("start local stream error: %w", err)
	}
	defer func() {
		_ = s.Stop(ctx)
	}()
	serverConn, err := t.Client.ConnectStream(ctx, streamName)
	if err != nil {
		return fmt.Errorf("connect to server error: %w", err)
	}
	if err := s.Join(ctx, serverConn); err != nil {
		return fmt.Errorf("join server connection to local stream error: %w", err)
	}

	// 设置输入流
	if f, ok := input.(*os.File); ok {
		// 将输入流设置为 raw 格式
		oldState, err := term.MakeRaw(int(f.Fd()))
		if err != nil {
			return fmt.Errorf("make input to raw error: %w", err)
		}
		defer func() {
			// 还原输入流
			_ = term.Restore(int(f.Fd()), oldState)
		}()
	}

	termConn := streams.NewPassthroughConnection(input, output)
	if err := s.Join(ctx, termConn); err != nil {
		return fmt.Errorf("join terminal io connection to local stream error: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case e, ok := <-s.ConnectionEvents():
			if !ok {
				return fmt.Errorf("stream closed")
			}
			if e.Type != streams.LeftEvent {
				logger.V(1).Info(fmt.Sprintf("received connection event: %s", e.Type))
				continue
			}
			if e.Connection == termConn {
				logger.V(1).Info("terminal closed")
				return nil
			} else {
				return fmt.Errorf("server connection closed")
			}
		}
	}
}
