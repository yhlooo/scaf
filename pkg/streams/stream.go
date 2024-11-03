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
	// ConnectionEvents 获取连接事件通道
	ConnectionEvents() <-chan ConnectionEvent
}

// ConnectionEvent 连接事件
type ConnectionEvent struct {
	// 事件类型
	Type ConnectionEventType
	// 连接
	Connection Connection
}

// ConnectionEventType 连接事件类型
type ConnectionEventType string

const (
	// JoinedEvent 连接已加入事件
	JoinedEvent ConnectionEventType = "Joined"
	// LeftEvent 连接已离开事件
	LeftEvent ConnectionEventType = "Left"
)
