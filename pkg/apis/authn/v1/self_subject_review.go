package v1

import (
	authnv1grpc "github.com/yhlooo/scaf/pkg/apis/authn/v1/grpc"
	metav1 "github.com/yhlooo/scaf/pkg/apis/meta/v1"
)

// SelfSubjectReview 检查自身身份
type SelfSubjectReview struct {
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Status SelfSubjectReviewStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// SelfSubjectReviewStatus 检查自身身份状态
type SelfSubjectReviewStatus struct {
	UserInfo UserInfo `json:"userInfo,omitempty" yaml:"userInfo,omitempty"`
}

// UserInfo 用户信息
type UserInfo struct {
	Username string `json:"username,omitempty" yaml:"username,omitempty"`
}

// NewSelfSubjectReviewFromGRPC 基于 *authnv1grpc.SelfSubjectReview 创建 *SelfSubjectReview
func NewSelfSubjectReviewFromGRPC(in *authnv1grpc.SelfSubjectReview) *SelfSubjectReview {
	if in == nil {
		return nil
	}
	meta := metav1.NewObjectMetaFromGRPC(in.GetMetadata())
	if meta == nil {
		meta = &metav1.ObjectMeta{}
	}
	return &SelfSubjectReview{
		ObjectMeta: *meta,
		Status: SelfSubjectReviewStatus{
			UserInfo: UserInfo{
				Username: in.GetStatus().GetUserInfo().GetUsername(),
			},
		},
	}
}

// NewGRPCSelfSubjectReview 基于 *SelfSubjectReview 创建 *authnv1grpc.SelfSubjectReview
func NewGRPCSelfSubjectReview(in *SelfSubjectReview) *authnv1grpc.SelfSubjectReview {
	if in == nil {
		return nil
	}
	return &authnv1grpc.SelfSubjectReview{
		Metadata: metav1.NewGRPCObjectMeta(&in.ObjectMeta),
		Status: &authnv1grpc.SelfSubjectReviewStatus{
			UserInfo: &authnv1grpc.UserInfo{
				Username: in.Status.UserInfo.Username,
			},
		},
	}
}
