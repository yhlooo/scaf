package apierrors

import (
	"net/http"

	metav1 "github.com/yhlooo/scaf/pkg/apis/meta/v1"
)

// 周知 metav1.Status Reason 的值
const (
	ReasonBadRequest          = "BadRequest"
	ReasonUnauthorized        = "Unauthorized"
	ReasonForbidden           = "Forbidden"
	ReasonNotFound            = "NotFound"
	ReasonInternalServerError = "InternalServerError"
)

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
