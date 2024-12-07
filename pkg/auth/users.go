package auth

import "github.com/yhlooo/scaf/pkg/randutil"

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

// IsStream 返回用户是否指定流
func IsStream(username string, streamName string) bool {
	return username == StreamUsername(streamName)
}

// StreamUsername 返回指定流用户名
func StreamUsername(streamName string) string {
	return StreamUsernamePrefix + streamName
}

// RandNormalUsername 生成一个随机用户名
func RandNormalUsername() string {
	return NormalUsernamePrefix + randutil.LowerAlphaNumeric(16)
}
