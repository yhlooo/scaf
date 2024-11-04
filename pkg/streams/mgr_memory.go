package streams

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
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

// CreateStream 创建并启动流
func (mgr *InMemoryManager) CreateStream(ctx context.Context, ins *StreamInstance) (*StreamInstance, error) {
	ins = ins.Clone()

	if err := ins.Stream.Start(ctx); err != nil {
		return nil, fmt.Errorf("start stream error: %w", err)
	}

	mgr.streamsLock.Lock()
	defer mgr.streamsLock.Unlock()
	ins.UID = UID(uuid.New().String())
	if mgr.streams == nil {
		mgr.streams = make(map[UID]*StreamInstance)
	}
	mgr.streams[ins.UID] = ins

	switch ins.StopPolicy {
	case OnFirstConnectionLeft:
		go stopStreamOnFirstConnectionLeft(ctx, ins)
	case OnBothConnectionsLeft:
		go stopStreamOnBothConnectionsLeft(ctx, ins)
	}

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

// DeleteStream 停止并删除流
func (mgr *InMemoryManager) DeleteStream(ctx context.Context, uid UID) error {
	mgr.streamsLock.RLock()
	stream, ok := mgr.streams[uid]
	mgr.streamsLock.RUnlock()
	if !ok {
		return fmt.Errorf("%w: stream %q not found", ErrStreamNotFound, uid)
	}

	// 停止流
	if err := stream.Stream.Stop(ctx); err != nil && !errors.Is(err, ErrStreamAlreadyStopped) {
		return fmt.Errorf("stop stream error: %w", err)
	}

	// 删除流
	mgr.streamsLock.Lock()
	delete(mgr.streams, uid)
	mgr.streamsLock.Unlock()

	return nil
}

// stopStreamOnFirstConnectionLeft 等待第一次连接离开时结束流
func stopStreamOnFirstConnectionLeft(ctx context.Context, ins *StreamInstance) {
	logger := logr.FromContextOrDiscard(ctx)
	for event := range ins.Stream.ConnectionEvents() {
		if event.Type == LeftEvent {
			if err := ins.Stream.Stop(ctx); err != nil {
				logger.Error(err, fmt.Sprintf("stop stream %q error", ins.UID))
			}
		}
	}
}

// stopStreamOnBothConnectionsLeft 等待所有连接都离开时结束流
func stopStreamOnBothConnectionsLeft(ctx context.Context, ins *StreamInstance) {
	logger := logr.FromContextOrDiscard(ctx)
	connCnt := 0
	for event := range ins.Stream.ConnectionEvents() {
		switch event.Type {
		case JoinedEvent:
			connCnt++
		case LeftEvent:
			connCnt--
		}
		if connCnt <= 0 {
			if err := ins.Stream.Stop(ctx); err != nil {
				logger.Error(err, fmt.Sprintf("stop stream %q error", ins.UID))
			}
		}
	}
}
