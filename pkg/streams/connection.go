package streams

import (
	"io"
)

// Connection 连接
type Connection interface {
	io.ReadWriteCloser
}
