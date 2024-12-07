package http

import (
	"net/http"
	"strings"

	"github.com/go-logr/logr"

	"github.com/yhlooo/scaf/pkg/randutil"
	"github.com/yhlooo/scaf/pkg/server/generic"
)

// GetTokenHandler 从请求中获取 Token 并注入上下文的 HTTP 处理器
func GetTokenHandler(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("Authorization")
		if token != "" && strings.HasPrefix(strings.ToLower(token), "bearer ") {
			token = token[7:]
			req = req.WithContext(generic.NewContextWithToken(req.Context(), token))
		}
		handler.ServeHTTP(w, req)
	}
}

// WithLoggerHandler 带 logr.Logger 的 HTTP 处理器
func WithLoggerHandler(handler http.Handler, logger logr.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		req = req.WithContext(logr.NewContext(req.Context(), logger.WithValues(
			"reqID", randutil.LowerAlphaNumeric(8),
			"request", req.Method+" "+req.RequestURI,
		)))
		handler.ServeHTTP(w, req)
	}
}
