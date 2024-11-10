package streams

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/creack/pty"
	"github.com/go-logr/logr"
)

var (
	resizeANSIRegexp = regexp.MustCompile("\x1b\\^(\\d+);(\\d+)s")
	ansiPMPrefix     = []byte("\x1b^")
)

// NewPTYConnection 创建 PTYConnection
func NewPTYConnection(name string, ptyIO *os.File) *PTYConnection {
	return &PTYConnection{
		f:      ptyIO,
		name:   name,
		logger: logr.Discard(),
	}
}

// PTYConnection 是 Connection 的基于伪终端的实现
type PTYConnection struct {
	f        *os.File
	name     string
	resizeCh chan *pty.Winsize

	logger logr.Logger
}

var _ Connection = &PTYConnection{}

// Name 返回连接名
func (conn *PTYConnection) Name() string {
	return conn.name
}

// Read 从 pty 输出流读
func (conn *PTYConnection) Read(p []byte) (int, error) {
	n, err := conn.f.Read(p)
	conn.logger.V(2).Info(fmt.Sprintf("read from pty: %q", p[:n]))
	return n, err
}

// Write 写到 pty 输入流
func (conn *PTYConnection) Write(p []byte) (int, error) {
	conn.logger.V(2).Info(fmt.Sprintf("write to pty: %q", p))

	if len(p) == 0 {
		return 0, nil
	}

	if !bytes.Contains(p, ansiPMPrefix) {
		// 没有 ANSI PM 序列
		return conn.f.Write(p)
	}

	ret := resizeANSIRegexp.FindAllSubmatch(p, -1)
	if len(ret) == 0 {
		// 没找到修改窗口大小的 ANSI 序列
		return conn.f.Write(p)
	}

	// PM<h>;<w>s
	sizeRe := ret[len(ret)-1]
	hStr := sizeRe[1]
	wStr := sizeRe[2]
	h, err := strconv.Atoi(string(hStr))
	if err != nil {
		return conn.f.Write(p)
	}
	w, err := strconv.Atoi(string(wStr))
	if err != nil {
		return conn.f.Write(p)
	}

	// 修改终端大小
	conn.logger.V(1).Info(fmt.Sprintf("set pty to size: h:%d w:%d", h, w))
	if err = pty.Setsize(conn.f, &pty.Winsize{
		Rows: uint16(h),
		Cols: uint16(w),
	}); err != nil {
		conn.logger.Error(err, "set pty size error")
	}

	pLen := len(p)
	p = resizeANSIRegexp.ReplaceAll(p, nil)
	err = nil
	if len(p) > 0 {
		_, err = conn.f.Write(p)
	}
	return pLen, err
}

// Close 关闭 pty 输入输出
func (conn *PTYConnection) Close() error {
	conn.logger.V(1).Info("close pty")
	return conn.f.Close()
}

// InjectLogger 注入 logger
func (conn *PTYConnection) InjectLogger(logger logr.Logger) {
	conn.logger = logger
}
