package streams

import (
	"context"
	"fmt"
	"sync"

	streamv1grpc "github.com/yhlooo/scaf/pkg/apis/stream/v1/grpc"
)

// NewGRPCStreamClientConnection 创建 GRPCStreamClientConnection
func NewGRPCStreamClientConnection(
	name string,
	client streamv1grpc.Streams_ConnectStreamClient,
) *GRPCStreamClientConnection {
	return &GRPCStreamClientConnection{
		name:   name,
		client: client,
	}
}

// GRPCStreamClientConnection 是 Connection 的基于 gRPC 流客户端的实现
type GRPCStreamClientConnection struct {
	name     string
	client   streamv1grpc.Streams_ConnectStreamClient
	closeErr error
}

var _ Connection = (*GRPCStreamClientConnection)(nil)

// Name 返回连接名
func (conn *GRPCStreamClientConnection) Name() string {
	return conn.name
}

// Send 发送
func (conn *GRPCStreamClientConnection) Send(_ context.Context, data []byte) error {
	if conn.closeErr != nil {
		return conn.closeErr
	}

	err := conn.client.Send(&streamv1grpc.Package{
		Content: data,
	})
	if err != nil {
		conn.closeErr = fmt.Errorf("%w: %s", ErrConnectionClosed, err.Error())
		_ = conn.client.CloseSend()
		return conn.closeErr
	}
	return nil
}

// Receive 接收
func (conn *GRPCStreamClientConnection) Receive(_ context.Context) ([]byte, error) {
	if conn.closeErr != nil {
		return nil, conn.closeErr
	}

	msg, err := conn.client.Recv()
	if err != nil {
		conn.closeErr = fmt.Errorf("%w: %s", ErrConnectionClosed, err.Error())
		_ = conn.client.CloseSend()
		return nil, conn.closeErr
	}
	return msg.GetContent(), nil
}

// Close 关闭连接
func (conn *GRPCStreamClientConnection) Close(_ context.Context) error {
	conn.closeErr = ErrConnectionClosed
	return conn.client.CloseSend()
}

// NewGRPCStreamServerConnection 创建 GRPCStreamServerConnection
func NewGRPCStreamServerConnection(
	name string,
	server streamv1grpc.Streams_ConnectStreamServer,
) *GRPCStreamServerConnection {
	return &GRPCStreamServerConnection{
		name:   name,
		server: server,
		done:   make(chan struct{}),
	}
}

// GRPCStreamServerConnection 是 Connection 的基于 gRPC 流服务端的实现
type GRPCStreamServerConnection struct {
	name      string
	server    streamv1grpc.Streams_ConnectStreamServer
	closeErr  error
	done      chan struct{}
	closeOnce sync.Once
}

var _ Connection = (*GRPCStreamServerConnection)(nil)

// Name 返回连接名
func (conn *GRPCStreamServerConnection) Name() string {
	return conn.name
}

// Send 发送
func (conn *GRPCStreamServerConnection) Send(_ context.Context, data []byte) error {
	if conn.closeErr != nil {
		return conn.closeErr
	}

	err := conn.server.Send(&streamv1grpc.Package{
		Content: data,
	})
	if err != nil {
		conn.closeErr = fmt.Errorf("%w: %s", ErrConnectionClosed, err.Error())
		return conn.closeErr
	}
	return nil
}

// Receive 接收
func (conn *GRPCStreamServerConnection) Receive(_ context.Context) ([]byte, error) {
	if conn.closeErr != nil {
		return nil, conn.closeErr
	}

	msg, err := conn.server.Recv()
	if err != nil {
		conn.closeErr = fmt.Errorf("%w: %s", ErrConnectionClosed, err.Error())
		return nil, conn.closeErr
	}
	return msg.GetContent(), nil
}

// Done 返回完成 channel
func (conn *GRPCStreamServerConnection) Done() <-chan struct{} {
	return conn.done
}

// Close 关闭连接
func (conn *GRPCStreamServerConnection) Close(_ context.Context) error {
	conn.closeOnce.Do(func() {
		conn.closeErr = ErrConnectionClosed
		if conn.done != nil {
			close(conn.done)
		}
	})
	return nil
}
