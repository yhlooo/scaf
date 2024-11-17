package grpc

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"

	"github.com/yhlooo/scaf/pkg/apierrors"
	metav1grpc "github.com/yhlooo/scaf/pkg/apis/meta/v1/grpc"
	streamv1 "github.com/yhlooo/scaf/pkg/apis/stream/v1"
	streamv1grpc "github.com/yhlooo/scaf/pkg/apis/stream/v1/grpc"
	"github.com/yhlooo/scaf/pkg/server/generic"
	"github.com/yhlooo/scaf/pkg/streams"
)

const (
	// MetadataKeyStreamName 表示流名的 metadata 键
	MetadataKeyStreamName = "scaf-stream-name"
	// MetadataKeyConnectionName 表示连接名的 metadata 键
	MetadataKeyConnectionName = "scaf-connection-name"
	// MetadataKeyToken 表示 Token 的 metadata 键
	MetadataKeyToken = "scaf-token"
)

// Options 选项
type Options struct {
	Logger logr.Logger
}

// NewStreamsServer 创建 gRPC 流服务
func NewStreamsServer(genericServer *generic.StreamsServer, opts Options) *StreamsServer {
	return &StreamsServer{
		genericServer: genericServer,
		logger:        opts.Logger,
	}
}

// StreamsServer 流服务
type StreamsServer struct {
	streamv1grpc.UnimplementedStreamsServer

	genericServer *generic.StreamsServer
	logger        logr.Logger
}

var _ streamv1grpc.StreamsServer = (*StreamsServer)(nil)

// CreateStream 创建流
func (s *StreamsServer) CreateStream(ctx context.Context, stream *streamv1grpc.Stream) (*streamv1grpc.Stream, error) {
	ctx = s.newContext(ctx, "request", "CreateStream")
	logger := logr.FromContextOrDiscard(ctx)
	logger.Info("request received")

	ret, err := s.genericServer.CreateStream(ctx, streamv1.NewStreamFromGRPC(stream))
	return streamv1.NewGRPCStream(ret), err
}

// GetStream 获取流
func (s *StreamsServer) GetStream(
	ctx context.Context,
	req *streamv1grpc.GetStreamRequest,
) (*streamv1grpc.Stream, error) {
	ctx = s.newContext(ctx, "request", "GetStream", "stream", req.GetName())
	logger := logr.FromContextOrDiscard(ctx)
	logger.Info("request received")

	ret, err := s.genericServer.GetStream(ctx, req.GetName())
	return streamv1.NewGRPCStream(ret), err
}

// ListStreams 列出流
func (s *StreamsServer) ListStreams(
	ctx context.Context,
	_ *streamv1grpc.ListStreamsRequest,
) (*streamv1grpc.StreamList, error) {
	ctx = s.newContext(ctx, "request", "ListStreams")
	logger := logr.FromContextOrDiscard(ctx)
	logger.Info("request received")

	ret, err := s.genericServer.ListStreams(ctx)
	return streamv1.NewGRPCStreamList(ret), err
}

// DeleteStream 删除流
func (s *StreamsServer) DeleteStream(
	ctx context.Context,
	req *streamv1grpc.DeleteStreamRequest,
) (*metav1grpc.Status, error) {
	ctx = s.newContext(ctx, "request", "DeleteStream", "stream", req.GetName())
	logger := logr.FromContextOrDiscard(ctx)
	logger.Info("request received")

	err := s.genericServer.DeleteStream(ctx, req.GetName())
	if err != nil {
		return nil, err
	}
	return &metav1grpc.Status{Code: 200, Reason: "Ok"}, nil
}

// ConnectStream 连接流
func (s *StreamsServer) ConnectStream(server streamv1grpc.Streams_ConnectStreamServer) error {
	ctx := server.Context()
	md, _ := metadata.FromIncomingContext(ctx)
	streamName := ""
	if values := md.Get(MetadataKeyStreamName); len(values) > 0 {
		streamName = values[0]
	}

	ctx = s.newContext(ctx, "request", "ConnectStream", "stream", streamName)
	logger := logr.FromContextOrDiscard(ctx)
	logger.Info("request received")

	ins, err := s.genericServer.GetStreamInstance(ctx, streamName)
	if err != nil {
		return err
	}

	connName := ""
	if values := md.Get(MetadataKeyConnectionName); len(values) > 0 {
		connName = values[0]
	}
	conn := streams.NewGRPCStreamServerConnection(connName, server)
	if err := ins.Stream.Join(ctx, conn); err != nil {
		logger.Error(err, "join stream error")
		return apierrors.NewInternalServerError(fmt.Errorf("join stream error: %w", err))
	}

	select {
	case <-ctx.Done():
	case <-conn.Done():
	}

	return nil
}

// newContext 创建请求上下文
func (s *StreamsServer) newContext(ctx context.Context, keyValues ...any) context.Context {
	// 注入 logger
	keyValues = append(keyValues, "reqID", uuid.New().String())
	ctx = logr.NewContext(ctx, s.logger.WithValues(keyValues...))

	// 注入 token
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
