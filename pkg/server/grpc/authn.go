package grpc

import (
	"context"

	"github.com/go-logr/logr"

	authnv1 "github.com/yhlooo/scaf/pkg/apis/authn/v1"
	authnv1grpc "github.com/yhlooo/scaf/pkg/apis/authn/v1/grpc"
	"github.com/yhlooo/scaf/pkg/server/generic"
)

// NewAuthenticationServer 创建 gRPC 认证服务
func NewAuthenticationServer(genericServer *generic.AuthenticationServer) *AuthenticationServer {
	return &AuthenticationServer{
		genericServer: genericServer,
	}
}

// AuthenticationServer 认证服务
type AuthenticationServer struct {
	authnv1grpc.UnimplementedAuthenticationServer
	genericServer *generic.AuthenticationServer
}

var _ authnv1grpc.AuthenticationServer = (*AuthenticationServer)(nil)

// CreateToken 创建 Token
func (s *AuthenticationServer) CreateToken(
	ctx context.Context,
	req *authnv1grpc.TokenRequest,
) (*authnv1grpc.TokenRequest, error) {
	logger := logr.FromContextOrDiscard(ctx)
	ctx = logr.NewContext(ctx, logger)
	logger.Info("request received")

	ret, err := s.genericServer.CreateToken(ctx, authnv1.NewTokenRequestFromGRPC(req))
	return authnv1.NewGRPCTokenRequest(ret), err
}

// CreateSelfSubjectReview 检查自身身份
func (s *AuthenticationServer) CreateSelfSubjectReview(
	ctx context.Context,
	review *authnv1grpc.SelfSubjectReview,
) (*authnv1grpc.SelfSubjectReview, error) {
	logger := logr.FromContextOrDiscard(ctx)
	ctx = logr.NewContext(ctx, logger)
	logger.Info("request received")

	ret, err := s.genericServer.CreateSelfSubjectReview(ctx, authnv1.NewSelfSubjectReviewFromGRPC(review))
	return authnv1.NewGRPCSelfSubjectReview(ret), err
}
