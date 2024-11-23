package streams

import (
	"context"

	metav1 "github.com/yhlooo/scaf/pkg/apis/meta/v1"
	streamv1 "github.com/yhlooo/scaf/pkg/apis/stream/v1"
)

// Manager 流管理器
type Manager interface {
	// CreateStream 创建流
	CreateStream(ctx context.Context, stream *StreamInstance) (*StreamInstance, error)
	// ListStreams 列出流
	ListStreams(ctx context.Context) ([]*StreamInstance, error)
	// GetStream 获取流
	GetStream(ctx context.Context, uid metav1.UID) (*StreamInstance, error)
	// DeleteStream 删除流
	DeleteStream(ctx context.Context, uid metav1.UID) error
}

// StreamInstance 流实例
type StreamInstance struct {
	Object streamv1.Stream
	Stream Stream
}

// Clone 返回流实例的一个拷贝
func (ins *StreamInstance) Clone() *StreamInstance {
	var annotations map[string]string
	if ins.Object.Annotations != nil {
		annotations = map[string]string{}
		for k, v := range ins.Object.Annotations {
			annotations[k] = v
		}
	}
	return &StreamInstance{
		Object: streamv1.Stream{
			ObjectMeta: metav1.ObjectMeta{
				Name:        ins.Object.ObjectMeta.Name,
				UID:         ins.Object.ObjectMeta.UID,
				Annotations: annotations,
			},
			Spec: streamv1.StreamSpec{
				StopPolicy: ins.Object.Spec.StopPolicy,
			},
			Status: streamv1.StreamStatus{
				Token: ins.Object.Status.Token,
			},
		},
		Stream: ins.Stream,
	}
}
