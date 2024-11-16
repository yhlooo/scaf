package v1

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	metav1grpc "github.com/yhlooo/scaf/pkg/apis/meta/v1/grpc"
)

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

// GRPCStatus 转为 gRPC 错误
func (s *Status) GRPCStatus() *status.Status {
	grpcCode := codes.Unknown
	switch s.Code {
	case http.StatusBadRequest:
		grpcCode = codes.InvalidArgument
	case http.StatusUnauthorized:
		grpcCode = codes.Unauthenticated
	case http.StatusForbidden:
		grpcCode = codes.PermissionDenied
	case http.StatusNotFound:
		grpcCode = codes.NotFound
	case http.StatusMethodNotAllowed:
		grpcCode = codes.Unimplemented
	case http.StatusInternalServerError:
		grpcCode = codes.Internal
	}
	ret, _ := status.New(grpcCode, s.Error()).WithDetails(&metav1grpc.Status{
		Code:    int32(s.Code),
		Reason:  s.Reason,
		Message: s.Message,
	})
	return ret
}
