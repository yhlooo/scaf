package v1

import (
	authnv1grpc "github.com/yhlooo/scaf/pkg/apis/authn/v1/grpc"
	metav1 "github.com/yhlooo/scaf/pkg/apis/meta/v1"
	metav1grpc "github.com/yhlooo/scaf/pkg/apis/meta/v1/grpc"
)

// TokenRequest Token 请求
type TokenRequest struct {
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Status TokenRequestStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// TokenRequestStatus Token 请求状态
type TokenRequestStatus struct {
	Token string `json:"token,omitempty" yaml:"token,omitempty"`
}

// NewTokenRequestFromGRPC 基于 *authnv1grpc.TokenRequest 创建 *TokenRequest
func NewTokenRequestFromGRPC(in *authnv1grpc.TokenRequest) *TokenRequest {
	if in == nil {
		return nil
	}
	return &TokenRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name:        in.GetMetadata().GetName(),
			UID:         metav1.UID(in.GetMetadata().GetUid()),
			Annotations: in.GetMetadata().GetAnnotations(),
		},
		Status: TokenRequestStatus{
			Token: in.GetStatus().GetToken(),
		},
	}
}

// NewGRPCTokenRequest 基于 *TokenRequest 创建 *authnv1grpc.TokenRequest
func NewGRPCTokenRequest(in *TokenRequest) *authnv1grpc.TokenRequest {
	if in == nil {
		return nil
	}
	return &authnv1grpc.TokenRequest{
		Metadata: &metav1grpc.ObjectMeta{
			Name:        in.Name,
			Uid:         string(in.UID),
			Annotations: in.Annotations,
		},
		Status: &authnv1grpc.TokenRequestStatus{
			Token: in.Status.Token,
		},
	}
}
