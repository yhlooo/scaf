package streams

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-logr/logr"
)

const (
	bufferedStreamLoggerName       = "buffered-stream"
	bufferedStreamBufferChannelLen = 256
	bufferedStreamRetryInterval    = time.Second
)

// NewBufferedStream 创建 BufferedStream
func NewBufferedStream() *BufferedStream {
	return &BufferedStream{
		eventCh: make(chan ConnectionEvent),
	}
}

// BufferedStream 带缓冲的流
type BufferedStream struct {
	lock   sync.RWMutex
	active bool

	connA       Connection
	connABuffCh chan []byte
	connB       Connection
	connBBuffCh chan []byte

	eventCh chan ConnectionEvent
}

var _ Stream = &BufferedStream{}

// Start 开始传输
func (s *BufferedStream) Start(_ context.Context) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.active {
		return ErrStreamAlreadyStarted
	}
	s.active = true
	return nil
}

// Join 将连接加入流
func (s *BufferedStream) Join(ctx context.Context, conn Connection) error {
	logger := logr.FromContextOrDiscard(ctx).WithName(bufferedStreamLoggerName)
	ctx = logr.NewContext(ctx, logger)

	if conn == nil {
		return fmt.Errorf("cannot join nil connection")
	}
	if logger.V(1).Enabled() {
		conn = ConnectionWithLog{Connection: conn}
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.active {
		return ErrStreamAlreadyStopped
	}

	switch {
	case s.connA == nil:
		s.connA = conn
		if s.connBBuffCh == nil {
			s.connBBuffCh = make(chan []byte, bufferedStreamBufferChannelLen)
		}
		go s.flushBuff(ctx, s.connABuffCh, s.connA)
		go s.handleConn(ctx, &s.connA, &s.connB, s.connBBuffCh)
	case s.connB == nil:
		s.connB = conn
		if s.connABuffCh == nil {
			s.connABuffCh = make(chan []byte, bufferedStreamBufferChannelLen)
		}
		go s.flushBuff(ctx, s.connBBuffCh, s.connB)
		go s.handleConn(ctx, &s.connB, &s.connA, s.connABuffCh)
	default:
		// 满员了，不能加入了
		return ErrStreamIsFull
	}

	select {
	case s.eventCh <- ConnectionEvent{Type: JoinedEvent, Connection: conn}:
	default:
	}

	return nil
}

// handleConn 处理连接
func (s *BufferedStream) handleConn(ctx context.Context, connRP, connWP *Connection, writeBuffCh chan<- []byte) {
	logger := logr.FromContextOrDiscard(ctx)

	s.lock.RLock()
	if *connRP == nil {
		s.lock.RUnlock()
		return
	}
	connR := *connRP
	s.lock.RUnlock()
	defer func() {
		_ = connR.Close(ctx)
		s.lock.Lock()
		*connRP = nil
		if s.eventCh != nil {
			select {
			case s.eventCh <- ConnectionEvent{Type: LeftEvent, Connection: connR}:
			default:
			}
		}
		s.lock.Unlock()
	}()

	for {
		data, err := connR.Receive(ctx)
		if err != nil {
			if errors.Is(err, ErrConnectionClosed) {
				// 连接已关闭
				return
			}
			logger.Error(err, "receive from connection error", "conn", connR.Name())
			time.Sleep(bufferedStreamRetryInterval)
			continue
		}

		s.lock.RLock()
		connW := *connWP
		s.lock.RUnlock()

		if connW == nil {
			// 另一个连接还未加入，先写到缓冲区
			select {
			case writeBuffCh <- bytes.Clone(data):
			default:
			}
			continue
		}

		// 另一个连接已经加入
		if err := connW.Send(ctx, data); err != nil {
			logger.Error(err, "send to connection error", "conn", connW.Name())
		}
	}
}

// flushBuff 将缓冲区的内容刷到连接
func (s *BufferedStream) flushBuff(ctx context.Context, buffCh <-chan []byte, conn Connection) {
	logger := logr.FromContextOrDiscard(ctx)

	for {
		select {
		case content, ok := <-buffCh:
			if !ok {
				return
			}
			if err := conn.Send(ctx, content); err != nil {
				logger.Error(err, "send to connection error", "conn", conn.Name())
			}
			logger.V(2).Info(fmt.Sprintf("send to connection: %q", content), "conn", conn.Name())
		default:
			return
		}
	}
}

// Stop 停止传输
func (s *BufferedStream) Stop(ctx context.Context) error {
	logger := logr.FromContextOrDiscard(ctx).WithName(bufferedStreamLoggerName)
	s.lock.Lock()
	defer s.lock.Unlock()
	if !s.active {
		return ErrStreamAlreadyStopped
	}

	if s.connA != nil {
		if err := s.connA.Close(ctx); err != nil {
			logger.Error(err, "close connection error", "conn", s.connA.Name())
		}
		s.connA = nil
	}
	if s.connABuffCh != nil {
		close(s.connABuffCh)
		s.connABuffCh = nil
	}
	if s.connB != nil {
		if err := s.connB.Close(ctx); err != nil {
			logger.Error(err, "close connection error", "conn", s.connB.Name())
		}
		s.connB = nil
	}
	if s.connBBuffCh != nil {
		close(s.connBBuffCh)
		s.connBBuffCh = nil
	}
	s.active = false
	close(s.eventCh)
	s.eventCh = nil

	return nil
}

// ConnectionEvents 获取连接事件通道
func (s *BufferedStream) ConnectionEvents() <-chan ConnectionEvent {
	return s.eventCh
}
