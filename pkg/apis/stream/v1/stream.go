package v1

import (
	metav1 "github.com/yhlooo/scaf/pkg/apis/meta/v1"
	metav1grpc "github.com/yhlooo/scaf/pkg/apis/meta/v1/grpc"
	streamv1grpc "github.com/yhlooo/scaf/pkg/apis/stream/v1/grpc"
)

// Stream 流
type Stream struct {
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Spec   StreamSpec   `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status StreamStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// StreamSpec 流定义
type StreamSpec struct {
	// 停止策略
	StopPolicy StreamStopPolicy `json:"stopPolicy,omitempty" yaml:"stopPolicy,omitempty"`
}

// StreamStopPolicy 流停止策略
type StreamStopPolicy string

const (
	// OnFirstConnectionLeft 第一次连接断开时停止
	OnFirstConnectionLeft StreamStopPolicy = "OnFirstConnectionLeft"
	// OnBothConnectionsLeft 两个连接都断开时停止
	OnBothConnectionsLeft StreamStopPolicy = "OnBothConnectionsLeft"
	// OnDelete 流被删除时停止
	OnDelete StreamStopPolicy = "OnDelete"
)

// StreamStatus 流状态
type StreamStatus struct {
	// 用于加入流的 token
	Token string `json:"token,omitempty" yaml:"token,omitempty"`
}

// StreamList 流列表
type StreamList struct {
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Items []Stream `json:"items,omitempty" yaml:"items,omitempty"`
}

// NewStreamFromGRPC 基于 *streamv1grpc.Stream 创建 *Stream
func NewStreamFromGRPC(in *streamv1grpc.Stream) *Stream {
	if in == nil {
		return nil
	}
	return &Stream{
		ObjectMeta: metav1.ObjectMeta{
			Name:        in.GetMetadata().GetName(),
			UID:         metav1.UID(in.GetMetadata().GetUid()),
			Annotations: in.GetMetadata().GetAnnotations(),
		},
		Spec: StreamSpec{
			StopPolicy: StreamStopPolicy(in.GetSpec().GetStopPolicy()),
		},
		Status: StreamStatus{
			Token: in.GetStatus().GetToken(),
		},
	}
}

// NewGRPCStream 基于 *Stream 创建 *streamv1grpc.Stream
func NewGRPCStream(in *Stream) *streamv1grpc.Stream {
	if in == nil {
		return nil
	}
	return &streamv1grpc.Stream{
		Metadata: &metav1grpc.ObjectMeta{
			Name:        in.Name,
			Uid:         string(in.UID),
			Annotations: in.Annotations,
		},
		Spec: &streamv1grpc.StreamSpec{
			StopPolicy: string(in.Spec.StopPolicy),
		},
		Status: &streamv1grpc.StreamStatus{
			Token: in.Status.Token,
		},
	}
}

// NewStreamListFromGRPC 基于 *streamv1grpc.StreamList 创建 *StreamList
func NewStreamListFromGRPC(in *streamv1grpc.StreamList) *StreamList {
	if in == nil {
		return nil
	}
	if len(in.Items) == 0 {
		return &StreamList{}
	}

	var items []Stream
	for _, item := range in.Items {
		items = append(items, *NewStreamFromGRPC(item))
	}
	return &StreamList{
		Items: items,
	}
}

// NewGRPCStreamList 基于 *StreamList 创建 *streamv1grpc.StreamList
func NewGRPCStreamList(in *StreamList) *streamv1grpc.StreamList {
	if in == nil {
		return nil
	}
	ret := &streamv1grpc.StreamList{
		Metadata: &metav1grpc.ListMeta{},
	}
	for _, item := range in.Items {
		ret.Items = append(ret.Items, NewGRPCStream(&item))
	}
	return ret
}
