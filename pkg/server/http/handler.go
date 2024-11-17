package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/yhlooo/scaf/pkg/apierrors"
	metav1 "github.com/yhlooo/scaf/pkg/apis/meta/v1"
	streamv1 "github.com/yhlooo/scaf/pkg/apis/stream/v1"
	"github.com/yhlooo/scaf/pkg/server/generic"
	"github.com/yhlooo/scaf/pkg/streams"
)

const (
	// ConnectionNameHeader 连接名头
	ConnectionNameHeader = "X-Scaf-Connection-Name"
)

// Options 选项
type Options struct {
	Logger logr.Logger
}

// NewHTTPHandler 创建 HTTP 请求处理器
func NewHTTPHandler(genericServer *generic.StreamsServer, opts Options) http.Handler {
	handlers := &httpHandlers{
		logger:        opts.Logger,
		genericServer: genericServer,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/streams", handlers.HandleCreateStream)
	mux.HandleFunc("GET /v1/streams", handlers.HandleListStreams)
	mux.HandleFunc("GET /v1/streams/{name}", handlers.HandleGetOrConnectStream)
	mux.HandleFunc("DELETE /v1/streams/{name}", handlers.HandleDeleteStream)
	return mux
}

// httpHandlers HTTP 请求处理器
type httpHandlers struct {
	genericServer *generic.StreamsServer
	logger        logr.Logger
}

// HandleCreateStream 处理创建流
func (h *httpHandlers) HandleCreateStream(w http.ResponseWriter, req *http.Request) {
	ctx := h.newContext(req, "request", "CreateStream")
	logger := logr.FromContextOrDiscard(ctx)
	logger.Info("request received")

	defer func() {
		_ = req.Body.Close()
	}()
	reqBody, err := io.ReadAll(io.LimitReader(req.Body, 1<<20))
	if err != nil {
		logger.Error(err, "read request error")
		responseStatus(ctx, w, apierrors.NewInternalServerError(fmt.Errorf("read request error: %w", err)))
		return
	}
	stream := &streamv1.Stream{}
	if err := json.Unmarshal(reqBody, stream); err != nil {
		logger.Error(err, "unmarshal request error")
		responseStatus(ctx, w, apierrors.NewBadRequestError(fmt.Errorf("parse request error: %w", err)))
		return
	}

	ret, err := h.genericServer.CreateStream(ctx, stream)
	if err != nil {
		responseStatus(ctx, w, apierrors.NewFromError(err))
		return
	}
	responseJSON(ctx, w, http.StatusCreated, ret)
}

// HandleListStreams 处理列出流
func (h *httpHandlers) HandleListStreams(w http.ResponseWriter, req *http.Request) {
	ctx := h.newContext(req, "request", "ListStreams")
	logger := logr.FromContextOrDiscard(ctx)
	logger.Info("request received")

	ret, err := h.genericServer.ListStreams(ctx)
	if err != nil {
		responseStatus(ctx, w, apierrors.NewFromError(err))
		return
	}
	responseJSON(ctx, w, http.StatusOK, ret)
}

// HandleGetOrConnectStream 处理获取或连接流
func (h *httpHandlers) HandleGetOrConnectStream(w http.ResponseWriter, req *http.Request) {
	streamName := req.PathValue("name")
	ctx := h.newContext(req, "request", "GetOrConnectStream", "stream", streamName)
	logger := logr.FromContextOrDiscard(ctx)
	logger.Info("request received")

	ins, err := h.genericServer.GetStreamInstance(ctx, streamName)
	if err != nil {
		responseStatus(ctx, w, apierrors.NewFromError(err))
		return
	}

	// 升级连接加入流
	if strings.ToLower(req.Header.Get("Connection")) == "upgrade" {
		connName := req.Header.Get(ConnectionNameHeader)
		switch {
		case websocket.IsWebSocketUpgrade(req):
			upgrader := &websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			}
			conn, err := upgrader.Upgrade(w, req, nil)
			if err != nil {
				logger.Error(err, "websocket upgrade error")
				responseStatus(ctx, w, apierrors.NewInternalServerError(
					fmt.Errorf("websocket upgrade error: %w", err),
				))
				return
			}
			if err := ins.Stream.Join(ctx, streams.NewWebSocketConnection(connName, conn)); err != nil {
				logger.Error(err, "join stream error")
				errMsg, _ := json.Marshal(apierrors.NewInternalServerError(fmt.Errorf("join stream error: %w", err)))
				if err := conn.WriteMessage(websocket.TextMessage, errMsg); err != nil {
					logger.Error(err, "send message error")
				}
				if err := conn.Close(); err != nil {
					logger.Error(err, "close websocket connection error")
				}
				return
			}
		default:
			// TODO: 支持其它协议
			responseStatus(ctx, w, apierrors.NewBadRequestError(
				fmt.Errorf("unsupported protocol: %s", req.Header.Get("Upgrade")),
			))
		}
		return
	}

	responseJSON(ctx, w, http.StatusOK, generic.NewStreamAPIObject(ins))
}

// HandleDeleteStream 处理删除流
func (h *httpHandlers) HandleDeleteStream(w http.ResponseWriter, req *http.Request) {
	streamName := req.PathValue("name")
	ctx := h.newContext(req, "request", "DeleteStream", "stream", streamName)
	logger := logr.FromContextOrDiscard(ctx)
	logger.Info("request received")

	err := h.genericServer.DeleteStream(ctx, req.PathValue("name"))
	if err != nil {
		responseStatus(ctx, w, apierrors.NewFromError(err))
		return
	}
	responseStatus(ctx, w, newOKStatus())
}

// newContext 创建请求上下文
func (h *httpHandlers) newContext(req *http.Request, keyValues ...any) context.Context {
	ctx := req.Context()

	// 注入 logger
	keyValues = append(keyValues, "reqID", uuid.New().String())
	ctx = logr.NewContext(ctx, h.logger.WithValues(keyValues...))

	// 注入 token
	token := req.Header.Get("Authorization")
	if token != "" && strings.HasPrefix(strings.ToLower(token), "bearer ") {
		token = token[7:]
		ctx = generic.NewContextWithToken(ctx, token)
	}

	return ctx
}

// newOKStatus 创建普通正常状态
func newOKStatus() *metav1.Status {
	return &metav1.Status{
		Code:   http.StatusOK,
		Reason: "OK",
	}
}

// responseJSON 发送 JSON 响应
func responseJSON(ctx context.Context, w http.ResponseWriter, code int, ret interface{}) {
	logger := logr.FromContextOrDiscard(ctx)

	raw, err := json.Marshal(ret)
	if err != nil {
		logger.Error(err, "marshal response to json error")
		code = http.StatusInternalServerError
		raw = []byte(`{"code":500,"reason":"InternalServer","message":"marshal response to json error"}`)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(raw); err != nil {
		logger.Error(err, "write response error")
	}
}

// responseStatus 发送 *metav1.Status 相应
func responseStatus(ctx context.Context, w http.ResponseWriter, status *metav1.Status) {
	responseJSON(ctx, w, status.Code, status)
}
