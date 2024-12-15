package bench

import (
	"context"
	"fmt"
	"hash/crc32"
	"sync"
	"time"

	"github.com/go-logr/logr"

	streamv1 "github.com/yhlooo/scaf/pkg/apis/stream/v1"
	"github.com/yhlooo/scaf/pkg/clients/common"
	"github.com/yhlooo/scaf/pkg/streams"
)

const (
	readWritePackageSize = 4 << 10 // 4KiB
)

// NewClient 创建基准测试客户端
func NewClient(client common.Client) *BenchmarkClient {
	return &BenchmarkClient{c: client}
}

// BenchmarkClient 基准测试客户端
type BenchmarkClient struct {
	c common.Client
}

// Report 测试报告
type Report struct {
	Ping      PingResult         `json:"ping,omitempty"`
	ReadOnly  TransmissionResult `json:"readOnly,omitempty"`
	WriteOnly TransmissionResult `json:"writeOnly,omitempty"`
	ReadWrite ReadWriteResult    `json:"readWrite,omitempty"`
}

// PingResult Ping 结果
type PingResult struct {
	// 往返时延
	RoundTripTime time.Duration `json:"roundTripTime"`
	// 丢包率
	LossRate float64 `json:"lossRate"`
}

// ReadWriteResult 读写结果
type ReadWriteResult struct {
	Read  TransmissionResult `json:"read,omitempty"`
	Write TransmissionResult `json:"write,omitempty"`
}

// TransmissionResult 读结果
type TransmissionResult struct {
	// 吞吐率（单位： Bytes/s ）
	Throughput uint64 `json:"throughput"`
	// 丢包率
	LossRate float64 `json:"lossRate"`
	// 成功传输的数据大小
	Size uint64 `json:"size"`
	// 成功传输的包数量
	Packages uint32 `json:"packages"`
}

// Run 运行基准测试
func (c *BenchmarkClient) Run(ctx context.Context, stream *streamv1.Stream) (*Report, error) {
	logger := logr.FromContextOrDiscard(ctx)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 与服务端建立连接
	conn, err := c.c.ConnectStream(ctx, stream.Name, common.ConnectStreamOptions{
		ConnectionName: "server",
	})
	if err != nil {
		return nil, fmt.Errorf("connect to server error: %w", err)
	}
	if logger.V(1).Enabled() {
		conn = streams.ConnectionWithLog{Connection: conn}
	}
	defer func() {
		cancel()
		_ = conn.Close(ctx)
	}()
	receiveMsgCh := make(chan Message)
	go c.runReceiveLoop(ctx, conn, receiveMsgCh)

	report := &Report{}

	// 测试往返时延
	pingRet, err := c.testPing(ctx, conn, receiveMsgCh, 5)
	if err != nil {
		return report, fmt.Errorf("ping error: %w", err)
	}
	report.Ping = *pingRet

	// 测试只读速率
	logger.Info("test read ...")
	readRet, _, err := c.testReadWrite(ctx, conn, receiveMsgCh, true, false, 10*time.Second)
	if err != nil {
		return report, fmt.Errorf("test read error: %w", err)
	}
	report.ReadOnly = *readRet
	logger.Info(fmt.Sprintf("received %d packages, size: %d", readRet.Packages, readRet.Size))

	logger.Info(fmt.Sprintf("waiting for connection idle ..."))
	if err := c.waitForConnIdle(ctx, receiveMsgCh, 5*time.Second); err != nil {
		return report, fmt.Errorf("wait for connection idle error: %w", err)
	}

	// 测试只写速率
	logger.Info("test write ...")
	_, writeRet, err := c.testReadWrite(ctx, conn, receiveMsgCh, false, true, 10*time.Second)
	if err != nil {
		return report, fmt.Errorf("test write error: %w", err)
	}
	report.WriteOnly = *writeRet
	logger.Info(fmt.Sprintf("sent %d packages, size: %d", writeRet.Packages, writeRet.Size))

	logger.Info(fmt.Sprintf("waiting for connection idle ..."))
	if err := c.waitForConnIdle(ctx, receiveMsgCh, 5*time.Second); err != nil {
		return report, fmt.Errorf("wait for connection idle error: %w", err)
	}

	// 测试同时读写速率
	logger.Info("test read and write ...")
	readRet, writeRet, err = c.testReadWrite(ctx, conn, receiveMsgCh, true, true, 10*time.Second)
	if err != nil {
		return report, fmt.Errorf("test read write error: %w", err)
	}
	report.ReadWrite.Read = *readRet
	report.ReadWrite.Write = *writeRet
	logger.Info(fmt.Sprintf("received %d packages, size: %d", readRet.Packages, readRet.Size))
	logger.Info(fmt.Sprintf("sent %d packages, size: %d", writeRet.Packages, writeRet.Size))

	return report, nil
}

// testPing 测试 Ping
func (c *BenchmarkClient) testPing(
	ctx context.Context,
	conn streams.Connection,
	msgCh <-chan Message,
	n int,
) (*PingResult, error) {
	logger := logr.FromContextOrDiscard(ctx)

	if n <= 0 {
		return &PingResult{}, nil
	}

	totalDuration := time.Duration(0)
	cnt := 0
mainLoop:
	for i := 0; i < n; i++ {
		if i != 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(500 * time.Millisecond):
			}
		}

		startTime := time.Now()
		if err := conn.Send(ctx, Ping(uint32(i)).Raw()); err != nil {
			logger.Error(err, fmt.Sprintf("send ping %d error", i))
			continue
		}
		for {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(5 * time.Second):
				logger.Info(fmt.Sprintf("WARN wait pong message %d timeout", i))
				continue mainLoop
			case msg, ok := <-msgCh:
				if !ok {
					return nil, fmt.Errorf("receive message channel closed")
				}
				switch typedMsg := msg.(type) {
				case Pong:
					if int(typedMsg) != i {
						logger.Info(fmt.Sprintf(
							"WARN received pong with wron seq: %d (expected: %d)",
							typedMsg, i,
						))
						continue
					}
				default:
					logger.Info(fmt.Sprintf("WARN received wrong message type: %s", msg.Type()))
					continue
				}
			}
			break
		}
		d := time.Since(startTime)
		logger.Info(fmt.Sprintf("ping %d: %s", i, d))
		totalDuration += d
		cnt++
	}
	return &PingResult{
		RoundTripTime: totalDuration / time.Duration(cnt),
		LossRate:      float64(n-cnt) / float64(n),
	}, nil
}

// testReadWrite 测试读写
func (c *BenchmarkClient) testReadWrite(
	ctx context.Context,
	conn streams.Connection,
	msgCh <-chan Message,
	read, write bool,
	d time.Duration,
) (readResult, writeResult *TransmissionResult, err error) {
	logger := logr.FromContextOrDiscard(ctx)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 发送开始指令
	mode := ReadWriteMode(0)
	if read {
		mode |= ReadMode
	}
	if write {
		mode |= WriteMode
	}
	if err := conn.Send(ctx, StartReadWrite{Mode: mode, ReadPackageSize: readWritePackageSize}.Raw()); err != nil {
		return nil, nil, fmt.Errorf("send start read write message error: %w", err)
	}

	var sendDone chan struct{}
	sendSeq := uint32(0)
	sendLock := sync.Mutex{}
	if write {
		// 发送数据
		sendDone = make(chan struct{})
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case <-sendDone:
					return
				default:
				}
				sendLock.Lock()
				sendSeq++
				if err := conn.Send(ctx, NewRandData(sendSeq, readWritePackageSize).Raw()); err != nil {
					logger.Error(err, fmt.Sprintf("send data %d error", sendSeq))
				}
				sendLock.Unlock()
			}
		}()
	}

	// 接收数据
	received := uint32(0)
	receivedSize := uint64(0)
	lastSeq := uint32(0)
	timer := time.NewTimer(d)
readDataLoop:
	for {
		var msg Message
		var ok bool
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		case <-timer.C:
			// 时间到了
			break readDataLoop
		case msg, ok = <-msgCh:
			if !ok {
				return nil, nil, fmt.Errorf("receive message channel closed")
			}
		}

		switch typedMsg := msg.(type) {
		case Data:
			if !read {
				// 不 read 应该不会收到数据类型的消息
				logger.Info(fmt.Sprintf("WARN received wrong message type: %s", msg.Type()))
				continue
			}
			// 校验
			if typedMsg.Seq() <= lastSeq {
				logger.Info(fmt.Sprintf("invalid data seq: %d (last: %d)", typedMsg.Seq(), lastSeq))
				continue
			}
			if sum := crc32.ChecksumIEEE(typedMsg.Content()); sum != typedMsg.Checksum() {
				logger.Info(fmt.Sprintf("invalid data checksum: %d (expected: %d)", sum, typedMsg.Checksum()))
				continue
			}
			received++
			receivedSize += uint64(len(typedMsg))
			lastSeq = typedMsg.Seq()
		default:
			logger.Info(fmt.Sprintf("WARN received wrong message type: %s", msg.Type()))
		}
	}

	// 发送结束指令
	sendLock.Lock()
	sendPkgs := sendSeq
	if err := conn.Send(ctx, StopReadWrite{}.Raw()); err != nil {
		logger.Error(err, fmt.Sprintf("send stop read write message error"))
	}
	sendLock.Unlock()

	readLossRate := float64(1)
	if lastSeq > 0 {
		readLossRate = float64(lastSeq-received) / float64(lastSeq)
	}
	readResult = &TransmissionResult{
		Throughput: receivedSize * 1000 / uint64(d.Milliseconds()),
		LossRate:   readLossRate,
		Size:       receivedSize,
		Packages:   received,
	}

	if write {
		// 等待写结果
		close(sendDone)
		timer = time.NewTimer(60 * time.Second)
		for {
			var msg Message
			var ok bool
			select {
			case <-ctx.Done():
				return readResult, nil, ctx.Err()
			case <-timer.C:
				return readResult, nil, fmt.Errorf("wait for write result timeout")
			case msg, ok = <-msgCh:
				if !ok {
					return readResult, nil, fmt.Errorf("receive message channel closed")
				}
			}
			switch typedMsg := msg.(type) {
			case Data:
				if !read {
					// 不 read 应该不会收到数据类型的消息
					logger.Info(fmt.Sprintf("WARN received wrong message type: %s", msg.Type()))
				}
				// 刚结束还有可能会收到一些数据，忽略
				continue
			case WriteResult:
				writeLossRate := float64(1)
				if sendPkgs > 0 {
					writeLossRate = float64(sendPkgs-typedMsg.ReceivedPackageCount) / float64(sendPkgs)
				}
				writeResult = &TransmissionResult{
					Throughput: uint64(typedMsg.ReceivedPackageCount) * readWritePackageSize * 1000 /
						uint64(d.Milliseconds()),
					LossRate: writeLossRate,
					Size:     uint64(typedMsg.ReceivedPackageCount) * readWritePackageSize,
					Packages: typedMsg.ReceivedPackageCount,
				}
				return readResult, writeResult, nil
			default:
				logger.Info(fmt.Sprintf("WARN received wrong message type: %s", msg.Type()))
			}
		}
	}

	return readResult, nil, nil
}

// waitForConnIdle 等待连接空闲
func (c *BenchmarkClient) waitForConnIdle(ctx context.Context, msgCh <-chan Message, idleTime time.Duration) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case _, ok := <-msgCh:
			if !ok {
				return fmt.Errorf("receive message channel closed")
			}
		case <-time.After(idleTime):
			return nil
		}
	}
}

// runReceiveLoop 运行接收循环
func (c *BenchmarkClient) runReceiveLoop(ctx context.Context, conn streams.Connection, msgCh chan<- Message) {
	defer close(msgCh)
	logger := logr.FromContextOrDiscard(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		raw, err := conn.Receive(ctx)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			logger.Error(err, "receive message error")
			continue
		}
		msg, err := ParseMessage(raw)
		if err != nil {
			logger.Error(err, "parse message error")
			continue
		}

		select {
		case <-ctx.Done():
			return
		case msgCh <- msg:
		}
	}
}
