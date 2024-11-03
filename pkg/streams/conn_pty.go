package streams

import (
	"bytes"
	"os"
	"regexp"
	"strconv"

	"github.com/creack/pty"
)

var (
	resizeANSIRegexp = regexp.MustCompile("\x1b\\^(\\d+);(\\d+)s")
	ansiPMPrefix     = []byte("\x1b^")
)

// NewPTYConnection 创建 PTYConnection
func NewPTYConnection(ptyIO *os.File) *PTYConnection {
	return &PTYConnection{
		File: ptyIO,
	}
}

// PTYConnection 是 Connection 的基于伪终端的实现
type PTYConnection struct {
	*os.File
	resizeCh chan *pty.Winsize
}

var _ Connection = &PTYConnection{}

// Write 写到 pty 输入流
func (conn *PTYConnection) Write(p []byte) (n int, err error) {
	if !bytes.Contains(p, ansiPMPrefix) {
		// 没有 ANSI PM 序列
		return conn.File.Write(p)
	}

	ret := resizeANSIRegexp.FindAllSubmatch(p, -1)
	if len(ret) == 0 {
		// 没找到修改窗口大小的 ANSI 序列
		return conn.File.Write(p)
	}

	// PM<h>;<w>s
	sizeRe := ret[len(ret)-1]
	hStr := sizeRe[1]
	wStr := sizeRe[2]
	h, err := strconv.Atoi(string(hStr))
	if err != nil {
		return conn.File.Write(p)
	}
	w, err := strconv.Atoi(string(wStr))
	if err != nil {
		return conn.File.Write(p)
	}

	// 修改终端大小
	_ = pty.Setsize(conn.File, &pty.Winsize{
		Rows: uint16(h),
		Cols: uint16(w),
	})

	return conn.File.Write(p)
}
