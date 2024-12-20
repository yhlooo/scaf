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

// StreamsServerOptions streams 服务选项
type StreamsServerOptions struct {
	TokenAuthenticator *auth.TokenAuthenticator
	StreamManager      streams.Manager
}

// NewStreamsServer 创建 *StreamsServer
func NewStreamsServer(opts StreamsServerOptions) *StreamsServer {
	return &StreamsServer{
		streamMgr:     opts.StreamManager,
		authenticator: opts.TokenAuthenticator,
	}
}

// StreamsServer 通用 streams 服务
type StreamsServer struct {
	streamMgr     streams.Manager
	authenticator *auth.TokenAuthenticator
}

// CreateStream 创建流
func (s *StreamsServer) CreateStream(ctx context.Context, stream *streamv1.Stream) (*streamv1.Stream, error) {
	logger := logr.FromContextOrDiscard(ctx)

	username, err := GetUsernameFromContext(ctx, s.authenticator)
	if err != nil {
		logger.Error(err, "get username error")
		return nil, apierrors.NewUnauthorizedError(err)
	}
	if !auth.IsAnonymous(username) && !auth.IsOwner(username, &stream.ObjectMeta) {
		stream.Owners = append(stream.Owners, username)
	}

	// 创建流
	strm := streams.NewBufferedStream()
	ins, err := s.streamMgr.CreateStream(ctx, &streams.StreamInstance{
		Object: *stream,
		Stream: strm,
	})
	if err != nil {
		logger.Error(err, "create stream error")
		return nil, apierrors.NewInternalServerError(err)
	}

	// 签发 token
	obj := &ins.Object
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
	return &ins.Object, nil
}

// GetStreamInstance 获取流实例
func (s *StreamsServer) GetStreamInstance(ctx context.Context, name string) (*streams.StreamInstance, error) {
	logger := logr.FromContextOrDiscard(ctx)

	username, err := GetUsernameFromContext(ctx, s.authenticator)
	if err != nil {
		logger.Error(err, "get username error")
		return nil, apierrors.NewUnauthorizedError(err)
	}
	if auth.IsAnonymous(username) {
		err := fmt.Errorf("user %q is not allowed to get stream %q", username, name)
		logger.Info(err.Error())
		return nil, apierrors.NewForbiddenError(err)
	}

	// 获取流
	ins, err := s.streamMgr.GetStream(ctx, metav1.UID(name))
	if err != nil {
		logger.Error(err, "get stream error")
		switch {
		case errors.Is(err, streams.ErrStreamNotFound):
			return nil, apierrors.NewNotFoundError(err)
		default:
			return nil, apierrors.NewInternalServerError(err)
		}
	}

	if !auth.IsStream(username, name) && !auth.IsOwner(username, &ins.Object.ObjectMeta) && !auth.IsAdmin(username) {
		err := fmt.Errorf("user %q is not allowed to get stream %q", username, name)
		logger.Info(err.Error())
		return nil, apierrors.NewForbiddenError(err)
	}

	return ins, nil
}

// ListStreams 列出流
func (s *StreamsServer) ListStreams(ctx context.Context) (*streamv1.StreamList, error) {
	logger := logr.FromContextOrDiscard(ctx)

	username, err := GetUsernameFromContext(ctx, s.authenticator)
	if err != nil {
		logger.Error(err, "get username error")
		return nil, apierrors.NewUnauthorizedError(err)
	}
	if auth.IsAnonymous(username) || auth.IsStreams(username) {
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
		if auth.IsOwner(username, &ins.Object.ObjectMeta) || auth.IsAdmin(username) {
			ret.Items = append(ret.Items, ins.Object)
		}
	}

	return ret, nil
}

// DeleteStream 删除流
func (s *StreamsServer) DeleteStream(ctx context.Context, name string) error {
	logger := logr.FromContextOrDiscard(ctx)

	username, err := GetUsernameFromContext(ctx, s.authenticator)
	if err != nil {
		logger.Error(err, "get username error")
		return apierrors.NewUnauthorizedError(err)
	}
	if auth.IsAnonymous(username) {
		err := fmt.Errorf("user %q is not allowed to delete stream %q", username, name)
		logger.Info(err.Error())
		return apierrors.NewForbiddenError(err)
	}
	if !auth.IsAdmin(username) && !auth.IsStream(username, name) {
		ins, err := s.streamMgr.GetStream(ctx, metav1.UID(name))
		if err != nil {
			logger.Error(err, "get stream error")
			switch {
			case errors.Is(err, streams.ErrStreamNotFound):
				return apierrors.NewNotFoundError(err)
			default:
				return apierrors.NewInternalServerError(err)
			}
		}
		if !auth.IsOwner(username, &ins.Object.ObjectMeta) {
			err := fmt.Errorf("user %q is not allowed to delete stream %q", username, name)
			logger.Info(err.Error())
			return apierrors.NewForbiddenError(err)
		}
	}

	if err := s.streamMgr.DeleteStream(ctx, metav1.UID(name)); err != nil {
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
