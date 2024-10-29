package streams

import (
	"io"
	"sync"

	"github.com/gorilla/websocket"
)

// NewWebSocketConnection 创建 WebSocketConnection
func NewWebSocketConnection(conn *websocket.Conn) *WebSocketConnection {
	return &WebSocketConnection{conn: conn}
}

// WebSocketConnection 是 Connection 的基于 WebSocket 的实现
type WebSocketConnection struct {
	conn   *websocket.Conn
	lock   sync.Mutex
	buf    []byte
	closed bool
}

var _ Connection = &WebSocketConnection{}

// Read 读数据
func (conn *WebSocketConnection) Read(p []byte) (int, error) {
	conn.lock.Lock()
	defer conn.lock.Unlock()

	if conn.closed {
		return 0, io.EOF
	}

	if len(conn.buf) == 0 {
		// 缓冲区没有了再读一点
		var err error
		_, conn.buf, err = conn.conn.ReadMessage()
		if err != nil {
			return 0, err
		}
	}

	n := len(p)
	if len(conn.buf) < n {
		n = len(conn.buf)
	}
	copy(p, conn.buf[:n])
	conn.buf = conn.buf[:n]

	return n, nil
}

// Write 写数据
func (conn *WebSocketConnection) Write(p []byte) (int, error) {
	conn.lock.Lock()
	defer conn.lock.Unlock()

	if conn.closed {
		return 0, io.EOF
	}

	if err := conn.conn.WriteMessage(websocket.BinaryMessage, p); err != nil {
		return 0, err
	}
	return len(p), nil
}

// Close 关闭连接
func (conn *WebSocketConnection) Close() error {
	conn.lock.Lock()
	defer conn.lock.Unlock()
	conn.closed = true
	return conn.conn.Close()
}
