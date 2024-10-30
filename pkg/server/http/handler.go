package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/gorilla/websocket"

	"github.com/yhlooo/scaf/pkg/apierrors"
	metav1 "github.com/yhlooo/scaf/pkg/apis/meta/v1"
	streamv1 "github.com/yhlooo/scaf/pkg/apis/stream/v1"
	"github.com/yhlooo/scaf/pkg/streams"
)

const (
	loggerName = "http"
)

// NewHTTPHandler 创建 HTTP 请求处理器
func NewHTTPHandler(ctx context.Context) http.Handler {
	logger := logr.FromContextOrDiscard(ctx).WithName(loggerName)

	handlers := &httpHandlers{
		logger:    logger,
		streamMgr: streams.NewInMemoryManager(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/streams", handlers.HandleCreateStream)
	mux.HandleFunc("GET /v1/streams", handlers.HandleListStreams)
	mux.HandleFunc("GET /v1/streams/{name}", handlers.HandleGetStream)
	mux.HandleFunc("DELETE /v1/streams/{name}", handlers.HandleDeleteStream)
	return mux
}

// httpHandlers HTTP 请求处理器
type httpHandlers struct {
	logger    logr.Logger
	streamMgr streams.Manager
}

// HandleCreateStream 处理创建流
func (h *httpHandlers) HandleCreateStream(w http.ResponseWriter, req *http.Request) {
	logger := h.logger.WithValues(
		"request", "CreateStream",
	)
	ctx := logr.NewContext(req.Context(), logger)
	req = req.WithContext(ctx)
	logger.Info("request received")

	s := streams.NewBufferedStream()
	ins, err := h.streamMgr.CreateStream(ctx, s)
	if err != nil {
		logger.Error(err, "create stream error")
		responseStatus(ctx, w, apierrors.NewInternalServerError(err))
		return
	}

	responseJSON(ctx, w, http.StatusCreated, newStreamAPIObject(ins))
}

// HandleListStreams 处理列出流
func (h *httpHandlers) HandleListStreams(w http.ResponseWriter, req *http.Request) {
	logger := h.logger.WithValues(
		"request", "ListStreams",
	)
	ctx := logr.NewContext(req.Context(), logger)
	req = req.WithContext(ctx)
	logger.Info("request received")

	streamList, err := h.streamMgr.ListStreams(ctx)
	if err != nil {
		logger.Error(err, "list stream error")
		responseStatus(ctx, w, apierrors.NewInternalServerError(err))
		return
	}

	ret := &streamv1.StreamList{
		TypeMeta: metav1.TypeMeta{
			APIVersion: streamv1.APIVersion,
			Kind:       streamv1.StreamListKind,
		},
	}
	for _, ins := range streamList {
		ret.Items = append(ret.Items, *newStreamAPIObject(ins))
	}

	responseJSON(ctx, w, http.StatusOK, ret)
}

// HandleGetStream 处理获取流
func (h *httpHandlers) HandleGetStream(w http.ResponseWriter, req *http.Request) {
	streamName := req.PathValue("name")
	logger := h.logger.WithValues(
		"request", "GetStream",
		"streamID", streamName,
	)
	ctx := logr.NewContext(req.Context(), logger)
	req = req.WithContext(ctx)
	logger.Info("request received")

	ins, err := h.streamMgr.GetStream(ctx, streams.UID(streamName))
	if err != nil {
		logger.Error(err, "get stream error")
		switch {
		case errors.Is(err, streams.ErrStreamNotFound):
			responseStatus(ctx, w, apierrors.NewNotFoundError(err))
		default:
			responseStatus(ctx, w, apierrors.NewInternalServerError(err))
		}
		return
	}

	// 升级连接加入流
	if req.Header.Get("Connection") == "upgrade" {
		switch {
		case websocket.IsWebSocketUpgrade(req):
			upgrader := &websocket.Upgrader{}
			conn, err := upgrader.Upgrade(w, req, nil)
			if err != nil {
				logger.Error(err, "upgrade websocket error")
				responseStatus(ctx, w, apierrors.NewInternalServerError(err))
				return
			}
			if err := ins.Stream.Join(ctx, streams.NewWebSocketConnection(conn)); err != nil {
				logger.Error(err, "join stream error")
				responseStatus(ctx, w, apierrors.NewInternalServerError(err))
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

	responseJSON(ctx, w, http.StatusOK, newStreamAPIObject(ins))
}

// HandleDeleteStream 处理删除流
func (h *httpHandlers) HandleDeleteStream(w http.ResponseWriter, req *http.Request) {
	streamName := req.PathValue("name")
	logger := h.logger.WithValues(
		"request", "DeleteStream",
		"streamID", streamName,
	)
	ctx := logr.NewContext(req.Context(), logger)
	req = req.WithContext(ctx)
	logger.Info("request received")

	if err := h.streamMgr.DeleteStream(ctx, streams.UID(streamName)); err != nil {
		logger.Error(err, "delete stream error")
		switch {
		case errors.Is(err, streams.ErrStreamNotFound):
			responseStatus(ctx, w, apierrors.NewNotFoundError(err))
		default:
			responseStatus(ctx, w, apierrors.NewInternalServerError(err))
		}
		return
	}
	responseStatus(ctx, w, newOKStatus())
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

// newStreamAPIObject 基于流实例创建流 API 对象
func newStreamAPIObject(ins *streams.StreamInstance) *streamv1.Stream {
	return &streamv1.Stream{
		TypeMeta: metav1.TypeMeta{
			APIVersion: streamv1.APIVersion,
			Kind:       streamv1.StreamKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: string(ins.UID), // TODO: 流名暂不支持自定义
			UID:  string(ins.UID),
		},
		Token: "", // TODO: ...
	}
}
