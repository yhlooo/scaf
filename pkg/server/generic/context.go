package generic

import (
	"context"

	"github.com/yhlooo/scaf/pkg/auth"
)

// tokenContextKey 上下文中存储 Token 的键
type tokenContextKey struct{}

// NewContextWithToken 创建包含 Token 的上下文
func NewContextWithToken(parent context.Context, token string) context.Context {
	return context.WithValue(parent, tokenContextKey{}, token)
}

// TokenFromContext 从上下文获取 Token
func TokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(tokenContextKey{}).(string)
	return token, ok
}

// GetUsernameFromContext 从上下文获取用户名
func GetUsernameFromContext(ctx context.Context, authenticator *auth.TokenAuthenticator) (string, error) {
	token, ok := TokenFromContext(ctx)
	if !ok || token == "" {
		return auth.AnonymousUsername, nil
	}
	return authenticator.AuthenticateToken(token)
}
