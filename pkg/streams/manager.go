package streams

import (
	"context"
)

// Manager 流管理器
type Manager interface {
	// CreateStream 创建流
	CreateStream(ctx context.Context, stream *StreamInstance) (*StreamInstance, error)
	// ListStreams 列出流
	ListStreams(ctx context.Context) ([]*StreamInstance, error)
	// GetStream 获取流
	GetStream(ctx context.Context, uid UID) (*StreamInstance, error)
	// DeleteStream 删除流
	DeleteStream(ctx context.Context, uid UID) error
}

// UID 流实例唯一 ID
type UID string

// StreamInstance 流实例
type StreamInstance struct {
	UID        UID
	StopPolicy StreamStopPolicy
	Stream     Stream
}

// StreamStopPolicy 流停止策略
type StreamStopPolicy string

const (
	// OnFirstConnectionLeft 第一次连接断开时停止
	OnFirstConnectionLeft StreamStopPolicy = "OnFirstConnectionLeft"
	// OnBothConnectionsLeft 两个连接都断开时停止
	OnBothConnectionsLeft StreamStopPolicy = "OnBothConnectionsLeft"
	// OnDelete 流被删除时停止
	OnDelete StreamStopPolicy = "OnDelete"
)

// Clone 返回流实例的一个拷贝
func (ins *StreamInstance) Clone() *StreamInstance {
	return &StreamInstance{
		UID:        ins.UID,
		StopPolicy: ins.StopPolicy,
		Stream:     ins.Stream,
	}
}
