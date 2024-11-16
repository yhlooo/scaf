package apierrors

import (
	"net/http"

	grpcstatus "google.golang.org/grpc/status"

	metav1 "github.com/yhlooo/scaf/pkg/apis/meta/v1"
)

// 周知 metav1.Status Reason 的值
const (
	ReasonUnknown             = "Unknown"
	ReasonBadRequest          = "BadRequest"
	ReasonUnauthorized        = "Unauthorized"
	ReasonForbidden           = "Forbidden"
	ReasonNotFound            = "NotFound"
	ReasonInternalServerError = "InternalServerError"
)

// NewFromError 从错误创建
//
//goland:noinspection GoTypeAssertionOnErrors
func NewFromError(err error) *metav1.Status {
	if status, ok := err.(*metav1.Status); ok {
		return status
	}

	if grpcStatus, ok := grpcstatus.FromError(err); ok {
		for _, d := range grpcStatus.Details() {
			if status, isStatus := d.(*metav1.Status); isStatus {
				return status
			}
		}
	}

	return &metav1.Status{
		Code:    http.StatusInternalServerError,
		Reason:  ReasonUnknown,
		Message: err.Error(),
	}
}

// NewBadRequestError 创建错误请求错误
func NewBadRequestError(err error) *metav1.Status {
	return &metav1.Status{
		Code:    http.StatusBadRequest,
		Reason:  ReasonBadRequest,
		Message: err.Error(),
	}
}

// NewUnauthorizedError 创建未认证请求错误
func NewUnauthorizedError(err error) *metav1.Status {
	return &metav1.Status{
		Code:    http.StatusUnauthorized,
		Reason:  ReasonUnauthorized,
		Message: err.Error(),
	}
}

// NewForbiddenError 创建不允许访问错误
func NewForbiddenError(err error) *metav1.Status {
	return &metav1.Status{
		Code:    http.StatusForbidden,
		Reason:  ReasonForbidden,
		Message: err.Error(),
	}
}

// NewNotFoundError 创建资源未找到结果
func NewNotFoundError(err error) *metav1.Status {
	return &metav1.Status{
		Code:    http.StatusNotFound,
		Reason:  ReasonNotFound,
		Message: err.Error(),
	}
}

// NewInternalServerError 创建服务内部错误结果
func NewInternalServerError(err error) *metav1.Status {
	return &metav1.Status{
		Code:    http.StatusInternalServerError,
		Reason:  ReasonInternalServerError,
		Message: err.Error(),
	}
}
