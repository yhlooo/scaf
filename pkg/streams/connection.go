package streams

import (
	"io"
)

// Connection 连接
type Connection interface {
	io.ReadWriteCloser
	// Name 返回连接名
	Name() string
}
