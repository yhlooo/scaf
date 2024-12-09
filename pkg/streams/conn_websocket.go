package streams

import (
	"context"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// NewWebSocketConnection 创建 WebSocketConnection
func NewWebSocketConnection(name string, conn *websocket.Conn) *WebSocketConnection {
	return &WebSocketConnection{
		name: name,
		conn: conn,
	}
}

// WebSocketConnection 是 Connection 的基于 WebSocket 的实现
type WebSocketConnection struct {
	name     string
	conn     *websocket.Conn
	closeErr error
	sendLock sync.Mutex
}

var _ Connection = &WebSocketConnection{}

// Name 返回连接名
func (conn *WebSocketConnection) Name() string {
	return conn.name
}

// Send 发送
func (conn *WebSocketConnection) Send(_ context.Context, data []byte) error {
	if conn.closeErr != nil {
		return conn.closeErr
	}

	conn.sendLock.Lock()
	defer conn.sendLock.Unlock()
	err := conn.conn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		if websocket.IsUnexpectedCloseError(err) {
			conn.closeErr = fmt.Errorf("%w: %s", ErrConnectionClosed, err.Error())
			return conn.closeErr
		}
		return err
	}

	return nil
}

// Receive 接收
func (conn *WebSocketConnection) Receive(_ context.Context) ([]byte, error) {
	if conn.closeErr != nil {
		return nil, conn.closeErr
	}

	_, msg, err := conn.conn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err) {
			conn.closeErr = fmt.Errorf("%w: %s", ErrConnectionClosed, err.Error())
			return nil, conn.closeErr
		}
		return nil, err
	}

	return msg, nil
}

// Close 关闭连接
func (conn *WebSocketConnection) Close(_ context.Context) error {
	conn.closeErr = ErrConnectionClosed
	return conn.conn.Close()
}
