package streams

import "context"

// Stream 流
type Stream interface {
	// Start 开始传输
	Start(ctx context.Context) error
	// Join 将连接加入流
	Join(ctx context.Context, conn Connection) error
	// Stop 停止传输
	Stop(ctx context.Context) error
}
