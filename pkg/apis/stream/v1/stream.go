package v1

import metav1 "github.com/yhlooo/scaf/pkg/apis/meta/v1"

const (
	// APIVersion API 版本
	APIVersion = "v1"
	// StreamKind Stream 类型名
	StreamKind = "Stream"
	// StreamListKind StreamList 类型名
	StreamListKind = "StreamList"
)

// Stream 流
type Stream struct {
	metav1.TypeMeta   `yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	// 用于加入流的 token
	Token string `json:"token,omitempty" yaml:"token,omitempty"`
}

// StreamList 流列表
type StreamList struct {
	metav1.TypeMeta `yaml:",inline"`
	Items           []Stream `json:"items,omitempty" yaml:"items,omitempty"`
}
