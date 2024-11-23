package exec

import (
	"encoding/json"
	"fmt"

	metav1 "github.com/yhlooo/scaf/pkg/apis/meta/v1"
	streamv1 "github.com/yhlooo/scaf/pkg/apis/stream/v1"
)

// 表示 exec 参数的注解
const (
	AnnoCommand      = "scaf/exec-command"
	AnnoInputEnabled = "scaf/exec-input-enabled"
	AnnoTTY          = "scaf/exec-tty"
)

// NewExecStream 创建 exec 流
func NewExecStream(command []string, input, tty bool) *streamv1.Stream {
	commandVal, _ := json.Marshal(command)
	inputVal := "false"
	if input {
		inputVal = "true"
	}
	ttyVal := "false"
	if tty {
		ttyVal = "true"
	}
	return &streamv1.Stream{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				AnnoCommand:      string(commandVal),
				AnnoInputEnabled: inputVal,
				AnnoTTY:          ttyVal,
			},
		},
		Spec: streamv1.StreamSpec{
			StopPolicy: streamv1.OnFirstConnectionLeft,
		},
	}
}

// GetExecOptions 通过流获取 exec 选项
func GetExecOptions(stream *streamv1.Stream) (command []string, input, tty bool, err error) {
	if stream == nil {
		return nil, false, false, fmt.Errorf("stream is nil")
	}
	err = json.Unmarshal([]byte(stream.Annotations[AnnoCommand]), &command)
	input = stream.Annotations[AnnoInputEnabled] == "true"
	tty = stream.Annotations[AnnoTTY] == "true"
	return
}
