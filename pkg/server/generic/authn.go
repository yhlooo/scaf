package generic

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/google/uuid"

	"github.com/yhlooo/scaf/pkg/apierrors"
	authnv1 "github.com/yhlooo/scaf/pkg/apis/authn/v1"
	metav1 "github.com/yhlooo/scaf/pkg/apis/meta/v1"
	"github.com/yhlooo/scaf/pkg/auth"
)

// AuthenticationServerOptions 认证服务选项
type AuthenticationServerOptions struct {
	TokenAuthenticator *auth.TokenAuthenticator
}

// NewAuthenticationServer 创建 *AuthenticationServer
func NewAuthenticationServer(opts AuthenticationServerOptions) *AuthenticationServer {
	return &AuthenticationServer{
		authenticator: opts.TokenAuthenticator,
	}
}

// AuthenticationServer 通用认证服务
type AuthenticationServer struct {
	authenticator *auth.TokenAuthenticator
}

// CreateToken 创建 Token
func (s *AuthenticationServer) CreateToken(ctx context.Context, _ *authnv1.TokenRequest) (*authnv1.TokenRequest, error) {
	logger := logr.FromContextOrDiscard(ctx)

	username := auth.RandNormalUsername()
	token, err := s.authenticator.IssueToken(username, 0)
	if err != nil {
		logger.Error(err, "issue token error")
		return nil, apierrors.NewInternalServerError(fmt.Errorf("issue token error: %w", err))
	}
	return &authnv1.TokenRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name: username,
			UID:  metav1.UID(uuid.New().String()),
		},
		Status: authnv1.TokenRequestStatus{
			Token: token,
		},
	}, nil
}

// CreateSelfSubjectReview 检查自身身份
func (s *AuthenticationServer) CreateSelfSubjectReview(ctx context.Context, _ *authnv1.SelfSubjectReview) (*authnv1.SelfSubjectReview, error) {
	logger := logr.FromContextOrDiscard(ctx)

	username, err := GetUsernameFromContext(ctx, s.authenticator)
	if err != nil {
		logger.Error(err, "get username error")
		return nil, apierrors.NewUnauthorizedError(err)
	}
	if auth.IsAnonymous(username) {
		return nil, apierrors.NewUnauthorizedError(fmt.Errorf("unauthorized"))
	}

	return &authnv1.SelfSubjectReview{
		Status: authnv1.SelfSubjectReviewStatus{
			UserInfo: authnv1.UserInfo{
				Username: username,
			},
		},
	}, nil
}
