package bench

import (
	"context"
	"fmt"
	"hash/crc32"

	"github.com/go-logr/logr"

	streamv1 "github.com/yhlooo/scaf/pkg/apis/stream/v1"
	"github.com/yhlooo/scaf/pkg/clients/common"
	"github.com/yhlooo/scaf/pkg/streams"
)

// NewServer 创建基准测试服务端
func NewServer(client common.Client) *BenchmarkServer {
	return &BenchmarkServer{c: client}
}

// BenchmarkServer 基准测试服务端
type BenchmarkServer struct {
	c common.Client
}

// Serve 运行测试服务
func (c *BenchmarkServer) Serve(ctx context.Context, stream *streamv1.Stream) error {
	logger := logr.FromContextOrDiscard(ctx)

	// 与服务端建立连接
	conn, err := c.c.ConnectStream(ctx, stream.Name, common.ConnectStreamOptions{
		ConnectionName: "server",
	})
	if err != nil {
		return fmt.Errorf("connect to server error: %w", err)
	}
	if logger.V(1).Enabled() {
		conn = streams.ConnectionWithLog{Connection: conn}
	}
	defer func() {
		_ = conn.Close(ctx)
	}()

	logger.Info("start serve benchmark")
	go c.handleBenchmarkRequests(ctx, conn)

	<-ctx.Done()
	return nil
}

// handleBenchmarkRequests 处理测试请求
func (c *BenchmarkServer) handleBenchmarkRequests(ctx context.Context, conn streams.Connection) {
	logger := logr.FromContextOrDiscard(ctx)

	for {
		raw, err := conn.Receive(ctx)
		if err != nil {
			if ctx.Err() == nil {
				logger.Error(err, "receive message error")
			}
			return
		}

		msg, err := ParseMessage(raw)
		if err != nil {
			logger.Error(err, "parse message error")
			continue
		}

		switch typedMsg := msg.(type) {
		case Ping:
			if err := conn.Send(ctx, Pong(typedMsg).Raw()); err != nil {
				logger.Error(err, fmt.Sprintf("send pong %d error", typedMsg))
			}
		case StartReadWrite:
			if err := c.handleReadWriteRequest(ctx, conn, typedMsg); err != nil {
				logger.Error(err, "handle read/write test request error")
			}
		default:
			logger.Info(fmt.Sprintf("ignore %s message", msg.Type()))
		}
	}
}

// handleReadWriteRequest 处理测试读写请求
func (c *BenchmarkServer) handleReadWriteRequest(
	ctx context.Context,
	conn streams.Connection,
	startMsg StartReadWrite,
) error {
	logger := logr.FromContextOrDiscard(ctx)

	read := startMsg.Mode&ReadMode != 0
	if read {
		pkgSize := startMsg.ReadPackageSize
		if pkgSize < 37 {
			pkgSize = 37
		}
		sendCTX, cancel := context.WithCancel(ctx)
		defer cancel()
		// 发送数据
		go func() {
			i := uint32(1)
			for {
				select {
				case <-sendCTX.Done():
					return
				default:
				}
				if err := conn.Send(ctx, NewRandData(i, pkgSize).Raw()); err != nil {
					logger.Error(err, fmt.Sprintf("send data %d error", i))
				}
				i++
			}
		}()
	}

	write := startMsg.Mode&WriteMode != 0

	// 接收数据
	received := uint32(0)
	lastSeq := uint32(0)
	for {
		raw, err := conn.Receive(ctx)
		if err != nil {
			return fmt.Errorf("receive message error: %w", err)
		}
		msg, err := ParseMessage(raw)
		if err != nil {
			logger.Error(err, "parse message error")
			continue
		}

		switch typedMsg := msg.(type) {
		case Data:
			if !write {
				// 不 write 应该不会收到数据类型的消息
				return fmt.Errorf("received wrong message type: %s", msg.Type())
			}
			// 校验
			if typedMsg.Seq <= lastSeq {
				logger.Info(fmt.Sprintf("invalid data seq: %d (last: %d)", typedMsg.Seq, lastSeq))
				continue
			}
			if sum := crc32.ChecksumIEEE(typedMsg.Content); sum != typedMsg.Checksum {
				logger.Info(fmt.Sprintf("invalid data checksum: %d (expected: %d)", sum, typedMsg.Checksum))
				continue
			}
			received++
			lastSeq = typedMsg.Seq
		case StopReadWrite:
			if write {
				// 返回成功接收的包数量
				if err := conn.Send(ctx, WriteResult{ReceivedPackageCount: received}.Raw()); err != nil {
					return fmt.Errorf("send write result message error: %w", err)
				}
			}
			return nil
		default:
			return fmt.Errorf("received wrong message type: %s", msg.Type())
		}
	}
}
