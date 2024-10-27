package http

import (
	"context"
	"net/http"

	"github.com/go-logr/logr"
)

const (
	loggerName = "http"
)

// NewHTTPHandler 创建 HTTP 请求处理器
func NewHTTPHandler(ctx context.Context) http.Handler {
	logger := logr.FromContextOrDiscard(ctx).WithName(loggerName)

	handlers := &httpHandlers{
		logger: logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/streams", handlers.HandleCreateStream)
	mux.HandleFunc("GET /v1/streams", handlers.HandleListStreams)
	mux.HandleFunc("GET /v1/streams/{id}", handlers.HandleGetStream)
	mux.HandleFunc("DELETE /v1/streams/{id}", handlers.HandleDeleteStream)
	return mux
}

type httpHandlers struct {
	logger logr.Logger
}

// HandleCreateStream 处理创建流
func (h *httpHandlers) HandleCreateStream(w http.ResponseWriter, req *http.Request) {
	logger := h.logger.WithValues(
		"request", "CreateStream",
	)
	req = req.WithContext(logr.NewContext(req.Context(), logger))
	logger.Info("request received")
	// TODO: ...
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("Ok"))
}

// HandleListStreams 处理列出流
func (h *httpHandlers) HandleListStreams(w http.ResponseWriter, req *http.Request) {
	logger := h.logger.WithValues(
		"request", "ListStreams",
	)
	req = req.WithContext(logr.NewContext(req.Context(), logger))
	logger.Info("request received")
	// TODO: ...
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Ok"))
}

// HandleGetStream 处理获取流
func (h *httpHandlers) HandleGetStream(w http.ResponseWriter, req *http.Request) {
	logger := h.logger.WithValues(
		"request", "GetStream",
		"streamID", req.PathValue("id"),
	)
	req = req.WithContext(logr.NewContext(req.Context(), logger))
	logger.Info("request received")
	// TODO: ...
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Ok"))
}

// HandleDeleteStream 处理删除流
func (h *httpHandlers) HandleDeleteStream(w http.ResponseWriter, req *http.Request) {
	logger := h.logger.WithValues(
		"request", "DeleteStream",
		"streamID", req.PathValue("id"),
	)
	req = req.WithContext(logr.NewContext(req.Context(), logger))
	logger.Info("request received")
	// TODO: ...
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Ok"))
}
