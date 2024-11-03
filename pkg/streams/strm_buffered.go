package streams

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/go-logr/logr"
)

const (
	bufferedStreamLoggerName       = "buffered-stream"
	bufferedStreamBufferChannelLen = 256     // 最大 256 * 4Ki = 1MBi
	bufferedStreamReadSize         = 4 << 10 // 4Ki
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
		return fmt.Errorf("already started")
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

	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.active {
		return fmt.Errorf("inactive stream")
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
		return fmt.Errorf("full stream")
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
		_ = connR.Close()
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

	tmp := make([]byte, bufferedStreamReadSize)
	for {
		n, err := connR.Read(tmp)
		if err != nil {
			if err == io.EOF {
				// 读完了
				return
			}
			logger.Error(err, "read from connection error")
			time.Sleep(bufferedStreamRetryInterval)
			continue
		}

		s.lock.RLock()
		connW := *connWP
		s.lock.RUnlock()

		if connW == nil {
			// 另一个连接还未加入，先写到缓冲区
			select {
			case writeBuffCh <- bytes.Clone(tmp[:n]):
			default:
			}
			continue
		}

		// 另一个连接已经加入
		if _, err := connW.Write(tmp[:n]); err != nil {
			logger.Error(err, "write to connection error")
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
			if _, err := conn.Write(content); err != nil {
				logger.Error(err, "write to connection error")
			}
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
		return fmt.Errorf("already stopped")
	}

	if s.connA != nil {
		if err := s.connA.Close(); err != nil {
			logger.Error(err, "close connection A error")
		}
		s.connA = nil
	}
	if s.connABuffCh != nil {
		close(s.connABuffCh)
		s.connABuffCh = nil
	}
	if s.connB != nil {
		if err := s.connB.Close(); err != nil {
			logger.Error(err, "close connection B error")
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
