package streams

import (
	"errors"
	"fmt"
	"io"
	"sync"
)

// NewPassthroughConnection 创建 PassthroughConnection
func NewPassthroughConnection(input io.Reader, output io.Writer) *PassthroughConnection {
	return &PassthroughConnection{
		input:  input,
		output: output,
	}
}

// PassthroughConnection 直通的 Connection
type PassthroughConnection struct {
	lock sync.RWMutex

	input  io.Reader
	output io.Writer

	closed bool
}

// Read 从输入读
func (conn *PassthroughConnection) Read(b []byte) (n int, err error) {
	conn.lock.RLock()
	defer conn.lock.RUnlock()
	if conn.closed {
		return 0, io.EOF
	}
	return conn.input.Read(b)
}

// Write 写到输出
func (conn *PassthroughConnection) Write(b []byte) (n int, err error) {
	conn.lock.RLock()
	defer conn.lock.RUnlock()
	if conn.closed {
		return 0, io.EOF
	}
	return conn.output.Write(b)
}

// Close 关闭连接
func (conn *PassthroughConnection) Close() error {
	var errs []error
	if closer, ok := conn.input.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close input error: %w", err))
		}
	}
	if closer, ok := conn.output.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close output error: %w", err))
		}
	}

	conn.lock.Lock()
	defer conn.lock.Unlock()

	conn.closed = true

	if len(errs) > 0 {
		return fmt.Errorf("close error: %w", errors.Join(errs...))
	}

	return nil
}
