package streams

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

// NewInMemoryManager 创建 InMemoryManager
func NewInMemoryManager() *InMemoryManager {
	return &InMemoryManager{}
}

// InMemoryManager 是 Manager 的基于内存的实现
type InMemoryManager struct {
	streamsLock sync.RWMutex
	streams     map[UID]*StreamInstance
}

var _ Manager = &InMemoryManager{}

// CreateStream 创建流
func (mgr *InMemoryManager) CreateStream(_ context.Context, stream Stream) (*StreamInstance, error) {
	mgr.streamsLock.Lock()
	defer mgr.streamsLock.Unlock()
	ins := NewSteamInstance(stream)
	if mgr.streams == nil {
		mgr.streams = make(map[UID]*StreamInstance)
	}
	mgr.streams[ins.UID] = ins
	return ins.Clone(), nil
}

// ListStreams 列出流
func (mgr *InMemoryManager) ListStreams(_ context.Context) ([]*StreamInstance, error) {
	mgr.streamsLock.RLock()
	defer mgr.streamsLock.RUnlock()
	if mgr.streams == nil {
		return nil, nil
	}
	streams := make([]*StreamInstance, 0, len(mgr.streams))
	for _, stream := range mgr.streams {
		streams = append(streams, stream.Clone())
	}
	// 按 uid 排序，保持结果稳定
	sort.Slice(streams, func(i, j int) bool {
		return streams[i].UID < streams[j].UID
	})
	return streams, nil
}

// GetStream 获取流
func (mgr *InMemoryManager) GetStream(_ context.Context, uid UID) (*StreamInstance, error) {
	mgr.streamsLock.RLock()
	defer mgr.streamsLock.RUnlock()
	stream, ok := mgr.streams[uid]
	if !ok {
		return nil, fmt.Errorf("%w: stream %q not found", ErrStreamNotFound, uid)
	}
	return stream.Clone(), nil
}

// DeleteStream 删除流
func (mgr *InMemoryManager) DeleteStream(_ context.Context, uid UID) error {
	mgr.streamsLock.Lock()
	defer mgr.streamsLock.Unlock()
	if _, ok := mgr.streams[uid]; !ok {
		return fmt.Errorf("%w: stream %q not found", ErrStreamNotFound, uid)
	}
	delete(mgr.streams, uid)
	return nil
}
