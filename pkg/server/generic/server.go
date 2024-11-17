package generic

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-logr/logr"

	"github.com/yhlooo/scaf/pkg/apierrors"
	metav1 "github.com/yhlooo/scaf/pkg/apis/meta/v1"
	streamv1 "github.com/yhlooo/scaf/pkg/apis/stream/v1"
	"github.com/yhlooo/scaf/pkg/auth"
	"github.com/yhlooo/scaf/pkg/streams"
)

// NewStreamsServer 创建 *StreamsServer
func NewStreamsServer(opts Options) *StreamsServer {
	return &StreamsServer{
		streamMgr:     opts.StreamManager,
		authenticator: opts.TokenAuthenticator,
	}
}

// Options 服务选项
type Options struct {
	TokenAuthenticator *auth.TokenAuthenticator
	StreamManager      streams.Manager
}

// StreamsServer 通用 streams 服务
type StreamsServer struct {
	streamMgr     streams.Manager
	authenticator *auth.TokenAuthenticator
}

// CreateStream 创建流
func (s *StreamsServer) CreateStream(ctx context.Context, stream *streamv1.Stream) (*streamv1.Stream, error) {
	logger := logr.FromContextOrDiscard(ctx)

	// 创建流
	strm := streams.NewBufferedStream()
	ins, err := s.streamMgr.CreateStream(ctx, &streams.StreamInstance{
		StopPolicy: streams.StreamStopPolicy(stream.Spec.StopPolicy),
		Stream:     strm,
	})
	if err != nil {
		logger.Error(err, "create stream error")
		return nil, apierrors.NewInternalServerError(err)
	}

	// 签发 token
	obj := NewStreamAPIObject(ins)
	obj.Status.Token, err = s.authenticator.IssueToken(auth.StreamUsername(obj.Name), 0)
	if err != nil {
		logger.Error(err, "issue stream token error")
		return nil, apierrors.NewInternalServerError(fmt.Errorf("issue stream token error: %w", err))
	}

	return obj, nil
}

// GetStream 获取流
func (s *StreamsServer) GetStream(ctx context.Context, name string) (*streamv1.Stream, error) {
	ins, err := s.GetStreamInstance(ctx, name)
	if err != nil {
		return nil, err
	}
	return NewStreamAPIObject(ins), nil
}

// GetStreamInstance 获取流实例
func (s *StreamsServer) GetStreamInstance(ctx context.Context, name string) (*streams.StreamInstance, error) {
	logger := logr.FromContextOrDiscard(ctx)

	username, err := s.getUsername(ctx)
	if err != nil {
		logger.Error(err, "get username error")
		return nil, apierrors.NewUnauthorizedError(err)
	}
	if !auth.IsAdmin(username) && !auth.IsStream(username, name) {
		err := fmt.Errorf("user %q is not allowed to get stream %q", username, name)
		logger.Info(err.Error())
		return nil, apierrors.NewForbiddenError(err)
	}

	// 获取流
	ins, err := s.streamMgr.GetStream(ctx, streams.UID(name))
	if err != nil {
		logger.Error(err, "get stream error")
		switch {
		case errors.Is(err, streams.ErrStreamNotFound):
			return nil, apierrors.NewNotFoundError(err)
		default:
			return nil, apierrors.NewInternalServerError(err)
		}
	}

	return ins, nil
}

// ListStreams 列出流
func (s *StreamsServer) ListStreams(ctx context.Context) (*streamv1.StreamList, error) {
	logger := logr.FromContextOrDiscard(ctx)

	username, err := s.getUsername(ctx)
	if err != nil {
		logger.Error(err, "get username error")
		return nil, apierrors.NewUnauthorizedError(err)
	}
	if !auth.IsAdmin(username) {
		err := fmt.Errorf("user %q is not allowed to list streams", username)
		logger.Info(err.Error())
		return nil, apierrors.NewForbiddenError(err)
	}

	streamList, err := s.streamMgr.ListStreams(ctx)
	if err != nil {
		logger.Error(err, "list stream error")
		return nil, apierrors.NewInternalServerError(err)
	}

	ret := &streamv1.StreamList{}
	for _, ins := range streamList {
		ret.Items = append(ret.Items, *NewStreamAPIObject(ins))
	}

	return ret, nil
}

// DeleteStream 删除流
func (s *StreamsServer) DeleteStream(ctx context.Context, name string) error {
	logger := logr.FromContextOrDiscard(ctx)

	username, err := s.getUsername(ctx)
	if err != nil {
		logger.Error(err, "get username error")
		return apierrors.NewUnauthorizedError(err)
	}
	if !auth.IsAdmin(username) && !auth.IsStream(username, name) {
		err := fmt.Errorf("user %q is not allowed to delete stream %q", username, name)
		logger.Info(err.Error())
		return apierrors.NewForbiddenError(err)
	}

	if err := s.streamMgr.DeleteStream(ctx, streams.UID(name)); err != nil {
		logger.Error(err, "delete stream error")
		switch {
		case errors.Is(err, streams.ErrStreamNotFound):
			return apierrors.NewNotFoundError(err)
		default:
			return apierrors.NewInternalServerError(err)
		}
	}

	return nil
}

// getUsername 获取请求用户名
func (s *StreamsServer) getUsername(ctx context.Context) (string, error) {
	token, ok := TokenFromContext(ctx)
	if !ok || token == "" {
		return auth.AnonymousUsername, nil
	}
	return s.authenticator.AuthenticateToken(token)
}

// NewStreamAPIObject 基于流实例创建流 API 对象
func NewStreamAPIObject(ins *streams.StreamInstance) *streamv1.Stream {
	return &streamv1.Stream{
		ObjectMeta: metav1.ObjectMeta{
			Name: string(ins.UID), // TODO: 流名暂不支持自定义
			UID:  string(ins.UID),
		},
		Spec: streamv1.StreamSpec{
			StopPolicy: streamv1.StreamStopPolicy(ins.StopPolicy),
		},
	}
}
