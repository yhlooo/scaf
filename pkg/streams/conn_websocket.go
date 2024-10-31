package streams

import (
	"bytes"
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
	conn         *websocket.Conn
	readBuffLock sync.Mutex
	readBuff     *bytes.Buffer
	closed       bool
}

var _ Connection = &WebSocketConnection{}

// Read 读数据
func (conn *WebSocketConnection) Read(p []byte) (int, error) {
	conn.readBuffLock.Lock()
	defer conn.readBuffLock.Unlock()

	if conn.closed {
		return 0, io.EOF
	}

	if conn.readBuff == nil {
		conn.readBuff = new(bytes.Buffer)
	}

	if conn.readBuff.Len() == 0 {
		// 缓冲区没有了再读一点
		var err error
		_, tmp, err := conn.conn.ReadMessage()
		if err != nil {
			_ = conn.Close()
			return 0, err
		}
		conn.readBuff.Write(tmp)
	}

	return conn.readBuff.Read(p)
}

// Write 写数据
func (conn *WebSocketConnection) Write(p []byte) (int, error) {
	if conn.closed {
		return 0, io.EOF
	}

	if err := conn.conn.WriteMessage(websocket.BinaryMessage, p); err != nil {
		_ = conn.Close()
		return 0, err
	}
	return len(p), nil
}

// Close 关闭连接
func (conn *WebSocketConnection) Close() error {
	conn.closed = true
	return conn.conn.Close()
}
