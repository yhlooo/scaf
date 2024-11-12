package streams

// Connection 连接
type Connection interface {
	// Name 返回连接名
	Name() string
	// Send 发送
	Send(data []byte) error
	// Receive 接收
	Receive() ([]byte, error)
	// Close 关闭连接
	Close() error
}
