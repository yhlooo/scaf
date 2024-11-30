package streams

import "context"

// Connection 连接
type Connection interface {
	// Name 返回连接名
	Name() string
	// Send 发送
	Send(ctx context.Context, data []byte) error
	// Receive 接收
	Receive(ctx context.Context) ([]byte, error)
	// Close 关闭连接
	Close(ctx context.Context) error
}
