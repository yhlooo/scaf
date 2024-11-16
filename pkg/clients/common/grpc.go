package common

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/yhlooo/scaf/pkg/apierrors"
	streamv1 "github.com/yhlooo/scaf/pkg/apis/stream/v1"
	streamv1grpc "github.com/yhlooo/scaf/pkg/apis/stream/v1/grpc"
	servergrpc "github.com/yhlooo/scaf/pkg/server/grpc"
	"github.com/yhlooo/scaf/pkg/streams"
)

const (
	defaultGRPCServerAddress = "localhost:9443"
)

// GRPCClientOptions gRPC 客户端选项
type GRPCClientOptions struct {
	// 服务端地址
	ServerAddress string
	// 用于认证的 Token
	Token string
}

// Complete 将选项补充完整
func (opts *GRPCClientOptions) Complete() {
	if opts.ServerAddress == "" {
		opts.ServerAddress = defaultGRPCServerAddress
	}
}

// NewGRPCClient 创建基于 gRPC 的客户端
func NewGRPCClient(opts GRPCClientOptions) (Client, error) {
	opts.Complete()
	conn, err := grpc.NewClient(opts.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &grpcClient{
		client: streamv1grpc.NewStreamsClient(conn),
	}, nil
}

// grpcClient 基于 gRPC 的客户端
type grpcClient struct {
	client streamv1grpc.StreamsClient
}

var _ Client = (*grpcClient)(nil)

// WithToken 返回使用指定 Token 的客户端
func (c *grpcClient) WithToken(token string) Client {
	//TODO: ...
	return c
}

// CreateStream 创建流
func (c *grpcClient) CreateStream(ctx context.Context, stream *streamv1.Stream) (*streamv1.Stream, error) {
	ret, err := c.client.CreateStream(ctx, streamv1.NewGRPCStream(stream))
	if err != nil {
		return nil, apierrors.NewFromError(err)
	}
	return streamv1.NewStreamFromGRPC(ret), nil
}

// GetStream 获取流
func (c *grpcClient) GetStream(ctx context.Context, name string) (*streamv1.Stream, error) {
	ret, err := c.client.GetStream(ctx, &streamv1grpc.GetStreamRequest{Name: name})
	if err != nil {
		return nil, apierrors.NewFromError(err)
	}
	return streamv1.NewStreamFromGRPC(ret), nil
}

// ListStreams 列出流
func (c *grpcClient) ListStreams(ctx context.Context) (*streamv1.StreamList, error) {
	ret, err := c.client.ListStreams(ctx, &streamv1grpc.ListStreamsRequest{})
	if err != nil {
		return nil, apierrors.NewFromError(err)
	}
	return streamv1.NewStreamListFromGRPC(ret), nil
}

// DeleteStream 删除流
func (c *grpcClient) DeleteStream(ctx context.Context, name string) error {
	_, err := c.client.DeleteStream(ctx, &streamv1grpc.DeleteStreamRequest{Name: name})
	if err != nil {
		return apierrors.NewFromError(err)
	}
	return nil
}

// ConnectStream 连接到流
func (c *grpcClient) ConnectStream(
	ctx context.Context,
	name string,
	opts ConnectStreamOptions,
) (streams.Connection, error) {
	md := metadata.MD{
		servergrpc.MetadataKeyStreamName:     []string{name},
		servergrpc.MetadataKeyConnectionName: []string{opts.ConnectionName},
	}
	ctx = metadata.NewOutgoingContext(ctx, md)
	streamClient, err := c.client.ConnectStream(ctx)
	if err != nil {
		return nil, apierrors.NewFromError(err)
	}
	return streams.NewGRPCStreamClientConnection(opts.ConnectionName, streamClient), nil
}
