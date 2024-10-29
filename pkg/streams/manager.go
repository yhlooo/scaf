package streams

import (
	"context"

	"github.com/google/uuid"
)

// Manager 流管理器
type Manager interface {
	// CreateStream 创建流
	CreateStream(ctx context.Context, stream Stream) (*StreamInstance, error)
	// ListStreams 列出流
	ListStreams(ctx context.Context) ([]*StreamInstance, error)
	// GetStream 获取流
	GetStream(ctx context.Context, uid UID) (*StreamInstance, error)
	// DeleteStream 删除流
	DeleteStream(ctx context.Context, uid UID) error
}

// UID 流实例唯一 ID
type UID string

// NewSteamInstance 创建流实例
func NewSteamInstance(stream Stream) *StreamInstance {
	return &StreamInstance{
		UID:    UID(uuid.New().String()),
		Stream: stream,
	}
}

// StreamInstance 流实例
type StreamInstance struct {
	UID    UID
	Stream Stream
}

// Clone 返回流实例的一个拷贝
func (ins *StreamInstance) Clone() *StreamInstance {
	return &StreamInstance{
		UID:    ins.UID,
		Stream: ins.Stream,
	}
}
