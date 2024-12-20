package bench

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"math/rand"
)

// ParseMessage 解析消息
func ParseMessage(raw []byte) (Message, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	switch raw[0] {
	case DataFlag:
		if len(raw) < 9 {
			return nil, fmt.Errorf("invalid data message length: %d (at lease 9 bytes)", len(raw))
		}
		return Data(raw), nil
	case PingFlag:
		if len(raw) != 5 {
			return nil, fmt.Errorf("invalid ping message: %v (must be 5 bytes)", raw)
		}
		return Ping(binary.BigEndian.Uint32(raw[1:])), nil
	case PongFlag:
		if len(raw) != 5 {
			return nil, fmt.Errorf("invalid pong message: %v (must be 5 bytes)", raw)
		}
		return Pong(binary.BigEndian.Uint32(raw[1:])), nil
	case StartReadWriteFlag:
		if len(raw) != 10 {
			return nil, fmt.Errorf("invalid start read write message: %v (must be 10 bytes)", raw)
		}
		return StartReadWrite{
			Mode:            ReadWriteMode(raw[1]),
			ReadPackageSize: binary.BigEndian.Uint64(raw[2:]),
		}, nil
	case StopReadWriteFlag:
		return StopReadWrite{}, nil
	case WriteResultFlag:
		if len(raw) != 5 {
			return nil, fmt.Errorf("invalid write result message: %v (must be 5 bytes)", raw)
		}
		return WriteResult{ReceivedPackageCount: binary.BigEndian.Uint32(raw[1:])}, nil
	default:
		return nil, fmt.Errorf("unknown data flag: %d, raw: %v", raw[0], raw)
	}
}

// Message 消息
type Message interface {
	// Type 返回消息类型
	Type() MessageType
	// Raw 返回消息原始数据
	Raw() []byte
}

// MessageType 消息类型
type MessageType string

const (
	DataType           MessageType = "Data"
	PingType           MessageType = "Ping"
	PongType           MessageType = "Pong"
	StartReadWriteType MessageType = "StartReadWrite"
	StopReadWriteType  MessageType = "StopReadWrite"
	WriteResultType    MessageType = "WriteResult"
)

const (
	DataFlag byte = iota
	PingFlag
	PongFlag
	StartReadWriteFlag
	StopReadWriteFlag
	WriteResultFlag
)

// NewRandData 创建随机数据
func NewRandData(seq uint32, size uint64) Data {
	if size < 9 {
		size = 9
	}
	data := make(Data, size)
	RenewRandData(data, seq)
	return data
}

// RenewRandData 重新生成随机数据
func RenewRandData(data Data, seq uint32) {
	content := data[9:]
	contentLen := len(content)
	randData := rand.Uint64()
	for i := 0; i < contentLen-7; i += 8 {
		binary.BigEndian.PutUint64(content[i:], randData)
	}
	data[0] = DataFlag
	binary.BigEndian.PutUint32(data[1:5], seq)
	binary.BigEndian.PutUint32(data[5:9], crc32.ChecksumIEEE(content))
}

// Data 数据消息
type Data []byte

var _ Message = Data(nil)

// Type 返回消息类型
func (d Data) Type() MessageType {
	return DataType
}

// Raw 返回消息原始数据
// 格式： 0(byte) seq(uint32) checksum(uint32) content([]byte)
func (d Data) Raw() []byte {
	return d
}

// Seq 返回包序号
func (d Data) Seq() uint32 {
	return binary.BigEndian.Uint32(d[1:5])
}

// Checksum 返回内容的 CRC32 校验和
func (d Data) Checksum() uint32 {
	return binary.BigEndian.Uint32(d[5:9])
}

// Content 返回内容
func (d Data) Content() []byte {
	return d[9:]
}

// Ping ping 消息
type Ping uint32

var _ Message = Ping(0)

// Type 返回消息类型
func (m Ping) Type() MessageType {
	return PingType
}

// Raw 返回消息原始数据
// 格式： 1(byte) useless([3]byte) seq(uint32)
func (m Ping) Raw() []byte {
	raw := make([]byte, 5)
	raw[0] = PingFlag
	binary.BigEndian.PutUint32(raw[1:], uint32(m))
	return raw
}

// Pong ping 消息
type Pong uint32

var _ Message = Pong(0)

// Type 返回消息类型
func (m Pong) Type() MessageType {
	return PongType
}

// Raw 返回消息原始数据
// 格式： 2(byte) useless([3]byte) seq(uint32)
func (m Pong) Raw() []byte {
	raw := make([]byte, 5)
	raw[0] = PongFlag
	binary.BigEndian.PutUint32(raw[1:], uint32(m))
	return raw
}

// ReadWriteMode 读写模式
type ReadWriteMode byte

const (
	ReadMode  ReadWriteMode = 1
	WriteMode ReadWriteMode = 2
)

// StartReadWrite 开始读写消息
type StartReadWrite struct {
	Mode            ReadWriteMode
	ReadPackageSize uint64
}

var _ Message = StartReadWrite{}

// Type 返回消息类型
func (m StartReadWrite) Type() MessageType {
	return StartReadWriteType
}

// Raw 返回消息原始数据
// 格式： 3(byte) mode(byte) readPackageSize(uint64)
func (m StartReadWrite) Raw() []byte {
	raw := make([]byte, 10)
	raw[0] = StartReadWriteFlag
	raw[1] = byte(m.Mode)
	binary.BigEndian.PutUint64(raw[2:], m.ReadPackageSize)
	return raw
}

// StopReadWrite 停止读写消息
type StopReadWrite struct{}

var _ Message = StopReadWrite{}

// Type 返回消息类型
func (m StopReadWrite) Type() MessageType {
	return StopReadWriteType
}

// Raw 返回消息原始数据
func (m StopReadWrite) Raw() []byte {
	return []byte{StopReadWriteFlag}
}

// WriteResult 写结果信息
type WriteResult struct {
	ReceivedPackageCount uint32
}

func (m WriteResult) Type() MessageType {
	return WriteResultType
}

// Raw 返回消息原始数据
// 格式： 5(byte) receivedPackageCount(uint32)
func (m WriteResult) Raw() []byte {
	raw := make([]byte, 5)
	raw[0] = WriteResultFlag
	binary.BigEndian.PutUint32(raw[1:], m.ReceivedPackageCount)
	return raw
}
