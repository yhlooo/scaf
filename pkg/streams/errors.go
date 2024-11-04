package streams

import "errors"

var (
	// ErrStreamAlreadyStarted 流已经开始了
	ErrStreamAlreadyStarted = errors.New("StreamAlreadyStarted")
	// ErrStreamNotFound 未找到流
	ErrStreamNotFound = errors.New("StreamNotFound")
	// ErrStreamIsFull 流满员了
	ErrStreamIsFull = errors.New("StreamIsFull")
	// ErrStreamAlreadyStopped 流已经停止了
	ErrStreamAlreadyStopped = errors.New("StreamAlreadyStopped")
)
