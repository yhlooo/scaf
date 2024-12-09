package grpc

import (
	"context"

	"github.com/go-logr/logr"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/yhlooo/scaf/pkg/server/generic"
	"github.com/yhlooo/scaf/pkg/utils/randutil"
)

const (
	// MetadataKeyStreamName 表示流名的 metadata 键
	MetadataKeyStreamName = "scaf-stream-name"
	// MetadataKeyConnectionName 表示连接名的 metadata 键
	MetadataKeyConnectionName = "scaf-connection-name"
	// MetadataKeyToken 表示 Token 的 metadata 键
	MetadataKeyToken = "scaf-token"
)

// WithLoggerInterceptor 往上下文注入 logr.Logger 的拦截器
func WithLoggerInterceptor(logger logr.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		return handler(logr.NewContext(ctx, logger.WithValues(
			"reqID", randutil.LowerAlphaNumeric(8),
			"method", info.FullMethod,
		)), req)
	}
}

// GetTokenInterceptor 获取 Token 的拦截器
func GetTokenInterceptor(
	ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	return handler(handleContextToken(ctx), req)
}

// WithLoggerStreamInterceptor 往上下文注入 logr.Logger 的拦截器
func WithLoggerStreamInterceptor(logger logr.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, &grpcmiddleware.WrappedServerStream{
			ServerStream: ss,
			WrappedContext: logr.NewContext(ss.Context(), logger.WithValues(
				"reqID", randutil.LowerAlphaNumeric(8),
				"method", info.FullMethod,
			)),
		})
	}
}

// GetTokenStreamInterceptor 获取 Token 的拦截器
func GetTokenStreamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	_ *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	return handler(srv, &grpcmiddleware.WrappedServerStream{
		ServerStream:   ss,
		WrappedContext: handleContextToken(ss.Context()),
	})
}

// handleContextToken 获取 metadata 中的 token 并注入 ctx
func handleContextToken(ctx context.Context) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	token := ""
	if values := md.Get(MetadataKeyToken); len(values) > 0 {
		token = values[0]
	}
	if token != "" {
		ctx = generic.NewContextWithToken(ctx, token)
	}
	return ctx
}
