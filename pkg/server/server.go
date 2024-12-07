package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/go-logr/logr"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"

	authnv1grpc "github.com/yhlooo/scaf/pkg/apis/authn/v1/grpc"
	streamv1grpc "github.com/yhlooo/scaf/pkg/apis/stream/v1/grpc"
	"github.com/yhlooo/scaf/pkg/auth"
	"github.com/yhlooo/scaf/pkg/server/generic"
	servergrpc "github.com/yhlooo/scaf/pkg/server/grpc"
	serverhttp "github.com/yhlooo/scaf/pkg/server/http"
	"github.com/yhlooo/scaf/pkg/streams"
)

const (
	loggerName        = "server"
	defaultListenAddr = ":9443"
)

// Options 是 Server 运行选项
type Options struct {
	// 监听地址
	ListenAddr string
	// Token 认证器选项
	TokenAuthenticator auth.TokenAuthenticatorOptions
}

// Complete 将选项补充完整
func (opts *Options) Complete() {
	if opts.ListenAddr == "" {
		opts.ListenAddr = defaultListenAddr
	}
}

// NewServer 创建 *Server
func NewServer(opts Options) *Server {
	opts.Complete()
	authenticator := auth.NewTokenAuthenticator(opts.TokenAuthenticator)
	streamMgr := streams.NewInMemoryManager()
	genericAuthnServer := generic.NewAuthenticationServer(generic.AuthenticationServerOptions{
		TokenAuthenticator: authenticator,
	})
	genericStreamsServer := generic.NewStreamsServer(generic.StreamsServerOptions{
		TokenAuthenticator: authenticator,
		StreamManager:      streamMgr,
	})
	return &Server{
		opts:                 opts,
		authenticator:        authenticator,
		streamMgr:            streamMgr,
		genericAuthnServer:   genericAuthnServer,
		genericStreamsServer: genericStreamsServer,
	}
}

// Server scaf 服务
type Server struct {
	opts Options

	startLock sync.RWMutex
	startOnce sync.Once
	cancel    context.CancelFunc
	done      chan struct{}

	listener net.Listener
	cmux     cmux.CMux

	httpListener net.Listener
	httpHandler  http.Handler

	grpcListener      net.Listener
	grpcServer        *grpc.Server
	grpcAuthnServer   *servergrpc.AuthenticationServer
	grpcStreamsServer *servergrpc.StreamsServer

	authenticator        *auth.TokenAuthenticator
	streamMgr            streams.Manager
	genericStreamsServer *generic.StreamsServer
	genericAuthnServer   *generic.AuthenticationServer
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

		// 监听端口
		s.listener, err = net.Listen("tcp", s.opts.ListenAddr)
		if err != nil {
			return
		}
		s.cmux = cmux.New(s.listener)
		// 根据协议分流
		s.grpcListener = s.cmux.MatchWithWriters(
			cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"),
		)
		s.httpListener = s.cmux.Match(cmux.Any())

		s.httpHandler = serverhttp.NewHTTPHandler(
			s.genericAuthnServer,
			s.genericStreamsServer,
			serverhttp.Options{
				Logger: logger.WithName("http"),
			},
		)

		s.grpcServer = grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				servergrpc.GetTokenInterceptor,
				servergrpc.WithLoggerInterceptor(logger.WithName("grpc")),
			),
			grpc.ChainStreamInterceptor(
				servergrpc.GetTokenStreamInterceptor,
				servergrpc.WithLoggerStreamInterceptor(logger.WithName("grpc")),
			),
		)
		s.grpcAuthnServer = servergrpc.NewAuthenticationServer(s.genericAuthnServer)
		authnv1grpc.RegisterAuthenticationServer(s.grpcServer, s.grpcAuthnServer)
		s.grpcStreamsServer = servergrpc.NewStreamsServer(s.genericStreamsServer)
		streamv1grpc.RegisterStreamsServer(s.grpcServer, s.grpcStreamsServer)

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

// Address 返回实际监听地址
func (s *Server) Address() net.Addr {
	s.startLock.RLock()
	defer s.startLock.RUnlock()

	if s.httpListener == nil {
		return nil
	}
	return s.httpListener.Addr()
}

// AdminToken 获取管理员用户 Token
func (s *Server) AdminToken() (string, error) {
	return s.authenticator.IssueToken(auth.AdminUsername, 0)
}

// run 运行服务，阻塞直到 ctx 被取消
func (s *Server) run(ctx context.Context) {
	logger := logr.FromContextOrDiscard(ctx)

	defer func() {
		s.grpcServer.GracefulStop()
		if err := s.listener.Close(); err != nil {
			logger.Error(err, "close tcp listener error")
		}
		close(s.done)
	}()

	cmuxDone := make(chan struct{})
	go func() {
		defer close(cmuxDone)
		if err := s.cmux.Serve(); err != nil {
			select {
			case <-ctx.Done():
				// ctx 结束了错误就没所谓了
				return
			default:
			}
			logger.Error(err, "cmux serve error")
		}
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

	grpcDone := make(chan struct{})
	go func() {
		defer close(grpcDone)
		if err := s.grpcServer.Serve(s.grpcListener); err != nil {
			select {
			case <-ctx.Done():
				// ctx 结束了错误就没所谓了
				return
			default:
			}
			logger.Error(err, "grpc serve error")
		}
	}()

	select {
	case <-ctx.Done():
	case <-cmuxDone:
	case <-httpDone:
	case <-grpcDone:
	}
}
