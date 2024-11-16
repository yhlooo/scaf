package v1

// ObjectMeta 对象元信息
type ObjectMeta struct {
	// 对象名
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// 对象全局唯一 ID
	UID string `json:"uid,omitempty" yaml:"uid,omitempty"`
}

// ListMeta 列表对象元信息
type ListMeta struct {
}
