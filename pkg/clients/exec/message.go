package exec

import (
	"encoding/binary"
	"fmt"
)

// ParseMessage 解析消息
func ParseMessage(raw []byte) (Message, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	switch raw[0] {
	case StdinDataFlag:
		return StdinData(raw[1:]), nil
	case StdoutDataFlag:
		return StdoutData(raw[1:]), nil
	case StderrDataFlag:
		return StderrData(raw[1:]), nil
	case ResizeFlag:
		if len(raw) != 5 {
			return nil, fmt.Errorf("invalid resize message: %v (must be 5 bytes)", raw)
		}
		return Resize{
			Height: binary.BigEndian.Uint16(raw[1:3]),
			Width:  binary.BigEndian.Uint16(raw[3:5]),
		}, nil
	case ExitCodeFlag:
		if len(raw) != 5 {
			return nil, fmt.Errorf("invalid exit code: %v (must be 5 bytes)", raw)
		}
		return ExitCode(int32(binary.BigEndian.Uint32(raw[1:5]))), nil
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
	StdinDataType  MessageType = "StdinData"
	StdoutDataType MessageType = "StdoutData"
	StderrDataType MessageType = "StderrData"
	ResizeType     MessageType = "Resize"
	ExitCodeType   MessageType = "ExitCode"
)

const (
	StdinDataFlag byte = iota
	StdoutDataFlag
	StderrDataFlag
	ResizeFlag
	ExitCodeFlag
)

// StdinData 标准输入流数据
type StdinData []byte

// Type 返回消息类型
func (d StdinData) Type() MessageType {
	return StdinDataType
}

// Raw 返回消息原始数据
// 格式： 0 data([]byte)
func (d StdinData) Raw() []byte {
	return append([]byte{StdinDataFlag}, d...)
}

// StdoutData 标准输出流数据
type StdoutData []byte

// Type 返回消息类型
func (d StdoutData) Type() MessageType {
	return StdoutDataType
}

// Raw 返回消息原始数据
// 格式： 1 data([]byte)
func (d StdoutData) Raw() []byte {
	return append([]byte{StdoutDataFlag}, d...)
}

// StderrData 标准错误流数据
type StderrData []byte

// Type 返回消息类型
func (d StderrData) Type() MessageType {
	return StderrDataType
}

// Raw 返回消息原始数据
// 格式： 2 data([]byte)
func (d StderrData) Raw() []byte {
	return append([]byte{StderrDataFlag}, d...)
}

// Resize 调整窗口大小消息
type Resize struct {
	Height uint16
	Width  uint16
}

// Type 返回消息类型
func (r Resize) Type() MessageType {
	return ResizeType
}

// Raw 返回消息原始数据
// 格式： 3 height(uint16) width(uint16)
func (r Resize) Raw() []byte {
	raw := make([]byte, 5)
	raw[0] = ResizeFlag
	binary.BigEndian.PutUint16(raw[1:3], r.Height)
	binary.BigEndian.PutUint16(raw[3:5], r.Width)
	return raw
}

// ExitCode 退出码
type ExitCode int32

func (e ExitCode) Type() MessageType {
	return ExitCodeType
}

// Raw 返回消息原始数据
// 格式： 4 code(int32)
func (e ExitCode) Raw() []byte {
	raw := make([]byte, 5)
	raw[0] = ExitCodeFlag
	binary.BigEndian.PutUint32(raw[1:5], uint32(e))
	return raw
}
