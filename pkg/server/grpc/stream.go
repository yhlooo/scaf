package grpc

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"google.golang.org/grpc/metadata"

	"github.com/yhlooo/scaf/pkg/apierrors"
	metav1grpc "github.com/yhlooo/scaf/pkg/apis/meta/v1/grpc"
	streamv1 "github.com/yhlooo/scaf/pkg/apis/stream/v1"
	streamv1grpc "github.com/yhlooo/scaf/pkg/apis/stream/v1/grpc"
	"github.com/yhlooo/scaf/pkg/server/generic"
	"github.com/yhlooo/scaf/pkg/streams"
)

// NewStreamsServer 创建 gRPC 流服务
func NewStreamsServer(genericServer *generic.StreamsServer) *StreamsServer {
	return &StreamsServer{
		genericServer: genericServer,
	}
}

// StreamsServer 流服务
type StreamsServer struct {
	streamv1grpc.UnimplementedStreamsServer
	genericServer *generic.StreamsServer
}

var _ streamv1grpc.StreamsServer = (*StreamsServer)(nil)

// CreateStream 创建流
func (s *StreamsServer) CreateStream(ctx context.Context, stream *streamv1grpc.Stream) (*streamv1grpc.Stream, error) {
	logger := logr.FromContextOrDiscard(ctx).WithValues("request", "CreateStream")
	ctx = logr.NewContext(ctx, logger)
	logger.Info("request received")

	ret, err := s.genericServer.CreateStream(ctx, streamv1.NewStreamFromGRPC(stream))
	return streamv1.NewGRPCStream(ret), err
}

// GetStream 获取流
func (s *StreamsServer) GetStream(
	ctx context.Context,
	req *streamv1grpc.GetStreamRequest,
) (*streamv1grpc.Stream, error) {
	logger := logr.FromContextOrDiscard(ctx).WithValues("request", "GetStream", "stream", req.GetName())
	ctx = logr.NewContext(ctx, logger)
	logger.Info("request received")

	ret, err := s.genericServer.GetStream(ctx, req.GetName())
	return streamv1.NewGRPCStream(ret), err
}

// ListStreams 列出流
func (s *StreamsServer) ListStreams(
	ctx context.Context,
	_ *streamv1grpc.ListStreamsRequest,
) (*streamv1grpc.StreamList, error) {
	logger := logr.FromContextOrDiscard(ctx).WithValues("request", "ListStreams")
	ctx = logr.NewContext(ctx, logger)
	logger.Info("request received")

	ret, err := s.genericServer.ListStreams(ctx)
	return streamv1.NewGRPCStreamList(ret), err
}

// DeleteStream 删除流
func (s *StreamsServer) DeleteStream(
	ctx context.Context,
	req *streamv1grpc.DeleteStreamRequest,
) (*metav1grpc.Status, error) {
	logger := logr.FromContextOrDiscard(ctx).WithValues("request", "DeleteStream", "stream", req.GetName())
	ctx = logr.NewContext(ctx, logger)
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

	logger := logr.FromContextOrDiscard(ctx).WithValues("request", "ConnectStream", "stream", streamName)
	ctx = logr.NewContext(ctx, logger)
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
