package auth

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	defaultSignKeyLength = 256
	nbfOffset            = -5 * time.Minute
)

// TokenAuthenticatorOptions TokenAuthenticator 选项
type TokenAuthenticatorOptions struct {
	Issuer  string
	SignKey []byte
}

// Complete 补全选项
func (opts *TokenAuthenticatorOptions) Complete() {
	// 生成密钥
	if len(opts.SignKey) == 0 {
		opts.SignKey = make([]byte, defaultSignKeyLength)
		_, _ = rand.Read(opts.SignKey)
	}
}

// NewTokenAuthenticator 创建 TokenAuthenticator
func NewTokenAuthenticator(opts TokenAuthenticatorOptions) *TokenAuthenticator {
	opts.Complete()
	a := &TokenAuthenticator{
		issuer: opts.Issuer,
		key:    make([]byte, len(opts.SignKey)),
	}
	copy(a.key, opts.SignKey)
	return a
}

// TokenAuthenticator 基于 Token 的认证器
type TokenAuthenticator struct {
	issuer string
	key    []byte
}

// AuthenticateToken 认证 Token
func (a *TokenAuthenticator) AuthenticateToken(token string) (username string, err error) {
	t, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		switch token.Method.Alg() {
		case jwt.SigningMethodHS256.Alg():
			return a.key, nil
		default:
			return nil, fmt.Errorf("unsupported signing method: %v", token.Method.Alg())
		}
	})
	if err != nil {
		return "", err
	}
	return t.Claims.GetSubject()
}

// IssueToken 签发 Token
func (a *TokenAuthenticator) IssueToken(username string, expire time.Duration) (token string, err error) {
	now := time.Now()
	var expiresAt *jwt.NumericDate
	if expire != 0 {
		expiresAt = jwt.NewNumericDate(now.Add(expire))
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    a.issuer,
		Subject:   username,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: expiresAt,
		NotBefore: jwt.NewNumericDate(now.Add(nbfOffset)),
	})
	return t.SignedString(a.key)
}
