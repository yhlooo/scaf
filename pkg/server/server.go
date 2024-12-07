package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/go-logr/logr"
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
	loggerName      = "server"
	defaultHTTPAddr = ":80"
	defaultGRPCAddr = ":9443"
)

// Options 是 Server 运行选项
type Options struct {
	// HTTP 监听地址
	HTTPAddr string
	// gRPC 监听地址
	GRPCAddr string
	// Token 认证器选项
	TokenAuthenticator auth.TokenAuthenticatorOptions
}

// Complete 将选项补充完整
func (opts *Options) Complete() {
	if opts.HTTPAddr == "" {
		opts.HTTPAddr = defaultHTTPAddr
	}
	if opts.GRPCAddr == "" {
		opts.GRPCAddr = defaultGRPCAddr
	}
}

// NewServer 创建 *Server
func NewServer(opts Options) *Server {
	opts.Complete()
	authenticator := auth.NewTokenAuthenticator(opts.TokenAuthenticator)
	streamMgr := streams.NewInMemoryManager()
	genericStreamsServer := generic.NewStreamsServer(generic.StreamsServerOptions{
		TokenAuthenticator: authenticator,
		StreamManager:      streamMgr,
	})
	genericAuthnServer := generic.NewAuthenticationServer(generic.AuthenticationServerOptions{
		TokenAuthenticator: authenticator,
	})
	return &Server{
		opts:                 opts,
		authenticator:        authenticator,
		streamMgr:            streamMgr,
		genericStreamsServer: genericStreamsServer,
		genericAuthnServer:   genericAuthnServer,
	}
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

	grpcListener      net.Listener
	grpcServer        *grpc.Server
	grpcStreamsServer *servergrpc.StreamsServer
	grpcAuthnServer   *servergrpc.AuthenticationServer

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

		s.httpListener, err = net.Listen("tcp", s.opts.HTTPAddr)
		if err != nil {
			return
		}
		s.httpHandler = serverhttp.NewHTTPHandler(
			s.genericStreamsServer,
			s.genericAuthnServer,
			serverhttp.Options{
				Logger: logger.WithName("http"),
			},
		)

		s.grpcListener, err = net.Listen("tcp", s.opts.GRPCAddr)
		if err != nil {
			return
		}
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
		s.grpcStreamsServer = servergrpc.NewStreamsServer(s.genericStreamsServer)
		streamv1grpc.RegisterStreamsServer(s.grpcServer, s.grpcStreamsServer)
		s.grpcAuthnServer = servergrpc.NewAuthenticationServer(s.genericAuthnServer)
		authnv1grpc.RegisterAuthenticationServer(s.grpcServer, s.grpcAuthnServer)

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

// AdminToken 获取管理员用户 Token
func (s *Server) AdminToken() (string, error) {
	return s.authenticator.IssueToken(auth.AdminUsername, 0)
}

// run 运行服务，阻塞直到 ctx 被取消
func (s *Server) run(ctx context.Context) {
	logger := logr.FromContextOrDiscard(ctx)

	defer func() {
		if err := s.httpListener.Close(); err != nil {
			logger.Error(err, "close http listener error")
		}
		s.grpcServer.GracefulStop()
		if err := s.grpcListener.Close(); err != nil {
			logger.Error(err, "close grpc listener error")
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
	case <-httpDone:
	case <-grpcDone:
	}
}
