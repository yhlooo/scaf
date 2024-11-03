package v1

import "fmt"

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

// Error 返回错误的字符串描述
func (s *Status) Error() string {
	return fmt.Sprintf("%s(%d): %s", s.Reason, s.Code, s.Message)
}
