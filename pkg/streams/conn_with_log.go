package streams

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/go-logr/logr"
)

// ConnectionWithLog 带日志的连接
type ConnectionWithLog struct {
	Connection
}

var _ Connection = &ConnectionWithLog{}

// Send 发送
func (conn ConnectionWithLog) Send(ctx context.Context, data []byte) error {
	logger := logr.FromContextOrDiscard(ctx).WithValues("conn", conn.Name())

	v := 1
	dataInfo := fmt.Sprintf("size: %d", len(data))
	if logger.V(2).Enabled() {
		v = 2
		dataInfo += fmt.Sprintf(", checksum: sha256:%x", sha256.Sum256(data))
	}

	err := conn.Connection.Send(ctx, data)
	if err != nil {
		logger.V(v).Info(fmt.Sprintf("send data error: %v, ", err) + dataInfo)
		return err
	}
	logger.V(v).Info("sent data, " + dataInfo)
	return nil
}

// Receive 接收
func (conn ConnectionWithLog) Receive(ctx context.Context) ([]byte, error) {
	logger := logr.FromContextOrDiscard(ctx).WithValues("conn", conn.Name())

	data, err := conn.Connection.Receive(ctx)
	if err != nil {
		logger.V(1).Info(fmt.Sprintf("receive error: %v", err))
		return data, err
	}

	if logger.V(2).Enabled() {
		logger.V(2).Info(fmt.Sprintf("received data, size: %d, checksum: sha256:%x", len(data), sha256.Sum256(data)))
	} else {
		logger.V(1).Info(fmt.Sprintf("received data, size: %d", len(data)))
	}

	return data, nil
}

// Close 关闭连接
func (conn ConnectionWithLog) Close(ctx context.Context) error {
	logger := logr.FromContextOrDiscard(ctx).WithValues("conn", conn.Name())
	err := conn.Connection.Close(ctx)
	if err != nil {
		logger.V(1).Info(fmt.Sprintf("close error: %v", err))
		return err
	}
	logger.V(1).Info("connection closed")
	return nil
}
