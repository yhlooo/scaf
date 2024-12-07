package v1

import metav1grpc "github.com/yhlooo/scaf/pkg/apis/meta/v1/grpc"

// ObjectMeta 对象元信息
type ObjectMeta struct {
	// 对象名
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// 对象全局唯一 ID
	UID UID `json:"uid,omitempty" yaml:"uid,omitempty"`
	// 注解
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	// 对象所有者用户名列表
	Owners []string `json:"owners,omitempty" yaml:"owners,omitempty"`
}

// UID 唯一 ID
type UID string

// ListMeta 列表对象元信息
type ListMeta struct {
}

// NewObjectMetaFromGRPC 基于 *metav1grpc.ObjectMeta 创建 *ObjectMeta
func NewObjectMetaFromGRPC(in *metav1grpc.ObjectMeta) *ObjectMeta {
	if in == nil {
		return nil
	}
	return &ObjectMeta{
		Name:        in.GetName(),
		UID:         UID(in.GetUid()),
		Annotations: in.GetAnnotations(),
		Owners:      in.GetOwners(),
	}
}

// NewGRPCObjectMeta 基于 *ObjectMeta 创建 *metav1grpc.ObjectMeta
func NewGRPCObjectMeta(in *ObjectMeta) *metav1grpc.ObjectMeta {
	if in == nil {
		return nil
	}
	return &metav1grpc.ObjectMeta{
		Name:        in.Name,
		Uid:         string(in.UID),
		Annotations: in.Annotations,
		Owners:      in.Owners,
	}
}
