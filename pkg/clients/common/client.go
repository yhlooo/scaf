package common

import (
	"context"

	streamv1 "github.com/yhlooo/scaf/pkg/apis/stream/v1"
	"github.com/yhlooo/scaf/pkg/streams"
)

// Client 客户端
type Client interface {
	// WithToken 返回使用指定 Token 的客户端
	WithToken(token string) Client
	// CreateStream 创建流
	CreateStream(ctx context.Context, stream *streamv1.Stream) (*streamv1.Stream, error)
	// GetStream 获取流
	GetStream(ctx context.Context, name string) (*streamv1.Stream, error)
	// ListStreams 列出流
	ListStreams(ctx context.Context) (*streamv1.StreamList, error)
	// DeleteStream 删除流
	DeleteStream(ctx context.Context, name string) error
	// ConnectStream 连接到流
	ConnectStream(ctx context.Context, name string, opts ConnectStreamOptions) (streams.Connection, error)
}

// ConnectStreamOptions 连接到流选项
type ConnectStreamOptions struct {
	ConnectionName string
}
