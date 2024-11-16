package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	"google.golang.org/grpc/metadata"

	"github.com/yhlooo/scaf/pkg/apierrors"
	metav1grpc "github.com/yhlooo/scaf/pkg/apis/meta/v1/grpc"
	streamv1grpc "github.com/yhlooo/scaf/pkg/apis/stream/v1/grpc"
	"github.com/yhlooo/scaf/pkg/auth"
	"github.com/yhlooo/scaf/pkg/streams"
)

const (
	loggerName = "grpc"

	// MetadataKeyStreamName 表示流名的 metadata 键
	MetadataKeyStreamName = "scaf-stream-name"
	// MetadataKeyConnectionName 表示连接名的 metadata 键
	MetadataKeyConnectionName = "scaf-connection-name"
)

// Options 选项
type Options struct {
	TokenAuthenticator *auth.TokenAuthenticator
	StreamManager      streams.Manager
}

// NewStreamsServer 创建 gRPC 流服务
func NewStreamsServer(ctx context.Context, opts Options) *StreamsServer {
	logger := logr.FromContextOrDiscard(ctx).WithName(loggerName)
	return &StreamsServer{
		logger:        logger,
		streamMgr:     opts.StreamManager,
		authenticator: opts.TokenAuthenticator,
	}
}

// StreamsServer 流服务
type StreamsServer struct {
	streamv1grpc.UnimplementedStreamsServer
	logger        logr.Logger
	streamMgr     streams.Manager
	authenticator *auth.TokenAuthenticator
}

var _ streamv1grpc.StreamsServer = (*StreamsServer)(nil)

// CreateStream 创建流
func (s *StreamsServer) CreateStream(ctx context.Context, stream *streamv1grpc.Stream) (*streamv1grpc.Stream, error) {
	logger := s.logger.WithValues(
		"request", "CreateStream",
	)
	ctx = logr.NewContext(ctx, logger)
	logger.Info("request received")

	strm := streams.NewBufferedStream()
	ins, err := s.streamMgr.CreateStream(ctx, &streams.StreamInstance{
		StopPolicy: streams.StreamStopPolicy(stream.GetSpec().GetStopPolicy()),
		Stream:     strm,
	})
	if err != nil {
		logger.Error(err, "create stream error")
		return nil, apierrors.NewInternalServerError(err)
	}

	obj := newStreamAPIObject(ins)
	obj.Status.Token, err = s.authenticator.IssueToken(auth.StreamUsername(obj.GetMetadata().GetName()), 0)
	if err != nil {
		logger.Error(err, "issue stream token error")
		return nil, apierrors.NewInternalServerError(fmt.Errorf("issue stream token error: %w", err))
	}

	return obj, nil
}

// GetStream 获取流
func (s *StreamsServer) GetStream(
	ctx context.Context,
	req *streamv1grpc.GetStreamRequest,
) (*streamv1grpc.Stream, error) {
	streamName := req.GetName()
	logger := s.logger.WithValues(
		"request", "GetStream",
		"streamID", streamName,
	)
	ctx = logr.NewContext(ctx, logger)
	logger.Info("request received")

	ins, err := s.streamMgr.GetStream(ctx, streams.UID(streamName))
	if err != nil {
		logger.Error(err, "get stream error")
		switch {
		case errors.Is(err, streams.ErrStreamNotFound):
			return nil, apierrors.NewNotFoundError(err)
		default:
			return nil, apierrors.NewInternalServerError(err)
		}
	}

	return newStreamAPIObject(ins), nil
}

// ListStreams 列出流
func (s *StreamsServer) ListStreams(
	ctx context.Context,
	_ *streamv1grpc.ListStreamsRequest,
) (*streamv1grpc.StreamList, error) {
	logger := s.logger.WithValues(
		"request", "ListStreams",
	)
	ctx = logr.NewContext(ctx, logger)
	logger.Info("request received")

	streamList, err := s.streamMgr.ListStreams(ctx)
	if err != nil {
		logger.Error(err, "list stream error")
		return nil, apierrors.NewInternalServerError(err)
	}

	ret := &streamv1grpc.StreamList{}
	for _, ins := range streamList {
		ret.Items = append(ret.Items, newStreamAPIObject(ins))
	}

	return ret, nil
}

// DeleteStream 删除流
func (s *StreamsServer) DeleteStream(
	ctx context.Context,
	req *streamv1grpc.DeleteStreamRequest,
) (*metav1grpc.Status, error) {
	streamName := req.GetName()
	logger := s.logger.WithValues(
		"request", "DeleteStream",
		"streamID", streamName,
	)
	ctx = logr.NewContext(ctx, logger)
	logger.Info("request received")

	if err := s.streamMgr.DeleteStream(ctx, streams.UID(streamName)); err != nil {
		logger.Error(err, "delete stream error")
		switch {
		case errors.Is(err, streams.ErrStreamNotFound):
			return nil, apierrors.NewNotFoundError(err)
		default:
			return nil, apierrors.NewInternalServerError(err)
		}
	}

	return &metav1grpc.Status{
		Code:   200,
		Reason: "Ok",
	}, nil
}

// ConnectStream 连接流
func (s *StreamsServer) ConnectStream(server streamv1grpc.Streams_ConnectStreamServer) error {
	ctx := server.Context()
	md, _ := metadata.FromIncomingContext(ctx)
	streamName := ""
	if values := md.Get(MetadataKeyStreamName); len(values) > 0 {
		streamName = values[0]
	}
	logger := s.logger.WithValues(
		"request", "ConnectStream",
		"streamID", streamName,
	)
	ctx = logr.NewContext(ctx, logger)
	logger.Info("request received")

	ins, err := s.streamMgr.GetStream(ctx, streams.UID(streamName))
	if err != nil {
		logger.Error(err, "get stream error")
		switch {
		case errors.Is(err, streams.ErrStreamNotFound):
			return apierrors.NewNotFoundError(err)
		default:
			return apierrors.NewInternalServerError(err)
		}
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

func newStreamAPIObject(ins *streams.StreamInstance) *streamv1grpc.Stream {
	return &streamv1grpc.Stream{
		Metadata: &metav1grpc.ObjectMeta{
			Name: string(ins.UID), // TODO: 流名暂不支持自定义
			Uid:  string(ins.UID),
		},
		Spec: &streamv1grpc.StreamSpec{
			StopPolicy: string(ins.StopPolicy),
		},
		Status: &streamv1grpc.StreamStatus{},
	}
}
