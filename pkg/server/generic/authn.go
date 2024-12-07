package generic

import (
	"context"
	"fmt"

	"github.com/google/uuid"

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
func (s *AuthenticationServer) CreateToken(_ context.Context, _ *authnv1.TokenRequest) (*authnv1.TokenRequest, error) {
	username := auth.RandNormalUsername()
	token, err := s.authenticator.IssueToken(username, 0)
	if err != nil {
		return nil, fmt.Errorf("issue token error: %w", err)
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
