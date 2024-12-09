package auth

import (
	"slices"
	"strings"

	metav1 "github.com/yhlooo/scaf/pkg/apis/meta/v1"
	"github.com/yhlooo/scaf/pkg/utils/randutil"
)

const (
	AdminUsername        = "system:admin"
	AnonymousUsername    = "system:anonymous"
	StreamUsernamePrefix = "system:stream:"
	NormalUsernamePrefix = "user:"
)

// IsAdmin 返回用户是否管理员
func IsAdmin(username string) bool {
	return username == AdminUsername
}

// IsAnonymous 返回用户是否未认证用户
func IsAnonymous(username string) bool {
	return username == AnonymousUsername
}

// IsStreams 返回用户是否是流
func IsStreams(username string) bool {
	return strings.HasPrefix(username, StreamUsernamePrefix)
}

// IsStream 返回用户是否指定流
func IsStream(username string, streamName string) bool {
	return username == StreamUsername(streamName)
}

// IsOwner 返回用户是否是指定对象所有者
func IsOwner(username string, meta *metav1.ObjectMeta) bool {
	if meta == nil {
		return false
	}
	return slices.Contains(meta.Owners, username)
}

// StreamUsername 返回指定流用户名
func StreamUsername(streamName string) string {
	return StreamUsernamePrefix + streamName
}

// RandNormalUsername 生成一个随机用户名
func RandNormalUsername() string {
	return NormalUsernamePrefix + randutil.LowerAlphaNumeric(16)
}
