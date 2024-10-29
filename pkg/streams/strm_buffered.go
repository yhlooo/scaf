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
	bufferedStreamLoggerName    = "buffered-stream"
	bufferedStreamMaxBufferLen  = 1 << 20 // 1MBi
	bufferedStreamReadSize      = 4 << 10 // 4Ki
	bufferedStreamRetryInterval = time.Second
)

// NewBufferedStream 创建 BufferedStream
func NewBufferedStream() *BufferedStream {
	return &BufferedStream{}
}

// BufferedStream 带缓冲的流
type BufferedStream struct {
	lock   sync.RWMutex
	active bool
	connA  Connection
	connB  Connection
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
		go s.handleConn(ctx, &s.connA, &s.connB)
	case s.connB == nil:
		s.connB = conn
		go s.handleConn(ctx, &s.connB, &s.connA)
	default:
		// 满员了，不能加入了
		return fmt.Errorf("full stream")
	}

	return nil
}

// handleConn 处理连接
func (s *BufferedStream) handleConn(ctx context.Context, connRP, connWP *Connection) {
	logger := logr.FromContextOrDiscard(ctx)

	s.lock.RLock()
	if *connRP == nil {
		return
	}
	connR := *connRP
	s.lock.RUnlock()
	defer func() {
		_ = connR.Close()
		s.lock.Lock()
		*connRP = nil
		s.lock.Unlock()
	}()

	buf := &bytes.Buffer{}
	for {
		tmp := make([]byte, bufferedStreamReadSize)
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
			// 另一个连接还未加入，先写到缓存
			if buf.Len()+n > bufferedStreamMaxBufferLen {
				// 缓冲区数据太多了，读出来一点
				ignore := make([]byte, bufferedStreamReadSize)
				if _, err = buf.Read(ignore); err != nil {
					logger.Error(err, "read from buffer error")
					time.Sleep(bufferedStreamRetryInterval)
					continue
				}
			}
			if _, err := buf.Write(tmp[:n]); err != nil {
				logger.Error(err, "write to buffer error")
				time.Sleep(bufferedStreamRetryInterval)
				continue
			}
			continue
		}

		// 另一个连接已经加入
		// 先发送缓冲区
		if buf.Len() > 0 {
			if _, err := io.Copy(connW, buf); err != nil {
				logger.Error(err, "copy from buffer to connection error")
			}
			buf.Reset()
		}
		// 然后发送读取的数据
		if _, err := connW.Write(tmp[:n]); err != nil {
			logger.Error(err, "write to connection error")
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
	if s.connB != nil {
		if err := s.connB.Close(); err != nil {
			logger.Error(err, "close connection B error")
		}
		s.connB = nil
	}
	s.active = false

	return nil
}
