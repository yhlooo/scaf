package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/go-logr/logr"

	serverhttp "github.com/yhlooo/scaf/pkg/server/http"
)

const (
	loggerName      = "server"
	defaultHTTPAddr = ":80"
)

// Options 是 Server 运行选项
type Options struct {
	// HTTP 监听地址
	HTTPAddr string
}

// Complete 将选项补充完整
func (opts *Options) Complete() {
	if opts.HTTPAddr == "" {
		opts.HTTPAddr = defaultHTTPAddr
	}
}

// NewServer 创建 *Server
func NewServer(opts Options) *Server {
	opts.Complete()
	return &Server{opts: opts}
}

// Server scaf 服务
type Server struct {
	opts Options

	startLock sync.RWMutex
	startOnce sync.Once
	cancel    context.CancelFunc
	done      chan struct{}

	httpListener net.Listener
	httpHandler  http.Handler
}

// Start 启动服务
// NOTE: 只能调用一次
func (s *Server) Start(ctx context.Context) error {
	s.startLock.Lock()
	defer s.startLock.Unlock()

	var err error
	start := false
	s.startOnce.Do(func() {
		start = true
		logger := logr.FromContextOrDiscard(ctx).WithName(loggerName)
		ctx, s.cancel = context.WithCancel(logr.NewContext(ctx, logger))
		s.done = make(chan struct{})

		s.httpListener, err = net.Listen("tcp", s.opts.HTTPAddr)
		if err != nil {
			return
		}
		s.httpHandler = serverhttp.NewHTTPHandler(ctx)

		go s.run(ctx)
	})

	if !start {
		return fmt.Errorf("already started")
	}

	return err
}

// Stop 停止服务
func (s *Server) Stop(ctx context.Context) error {
	s.startLock.RLock()
	defer s.startLock.RUnlock()

	if s.cancel != nil {
		s.cancel()
	}

	if s.done == nil {
		return nil
	}

	// 等待服务结束
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.done:
	}

	return nil
}

// Done 返回服务 done channel ，当服务结束运行时该 channel 会被关闭
// NOTE: 必须在 Start 方法被调用后调用，否则返回 nil
func (s *Server) Done() <-chan struct{} {
	s.startLock.RLock()
	defer s.startLock.RUnlock()

	return s.done
}

// HTTPAddr 返回 HTTP 实际监听地址
func (s *Server) HTTPAddr() net.Addr {
	s.startLock.RLock()
	defer s.startLock.RUnlock()

	if s.httpListener == nil {
		return nil
	}
	return s.httpListener.Addr()
}

// run 运行服务，阻塞直到 ctx 被取消
func (s *Server) run(ctx context.Context) {
	logger := logr.FromContextOrDiscard(ctx)

	defer func() {
		if err := s.httpListener.Close(); err != nil {
			logger.Error(err, "close http listener error")
		}
		close(s.done)
	}()

	httpDone := make(chan struct{})
	go func() {
		defer close(httpDone)
		if err := http.Serve(s.httpListener, s.httpHandler); err != nil {
			select {
			case <-ctx.Done():
				// ctx 结束了错误就没所谓了
				return
			default:
			}
			logger.Error(err, "http serve error")
		}
	}()

	select {
	case <-ctx.Done():
	case <-httpDone:
	}
}
