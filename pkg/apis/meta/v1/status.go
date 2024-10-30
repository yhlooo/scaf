package v1

// TypeMeta 数据类型元信息
type TypeMeta struct {
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`
}

// ObjectMeta 对象元信息
type ObjectMeta struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	UID  string `json:"uid,omitempty" yaml:"uid,omitempty"`
}

// Status 结果状态
type Status struct {
	Code    int    `json:"code" yaml:"code"`
	Reason  string `json:"reason,omitempty" yaml:"reason,omitempty"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
}
