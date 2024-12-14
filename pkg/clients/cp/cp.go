package cp

import (
	"archive/tar"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-logr/logr"

	streamv1 "github.com/yhlooo/scaf/pkg/apis/stream/v1"
	"github.com/yhlooo/scaf/pkg/clients/common"
	"github.com/yhlooo/scaf/pkg/streams"
)

const (
	maxReadSize  = 16 << 10
	startRecvMsg = "StartReceiving"
	sendDoneMsg  = "SendCompleted"
	recvDoneMsg  = "ReceivingCompleted"
)

// New 创建 CopyFileClient
func New(client common.Client) *CopyFileClient {
	return &CopyFileClient{c: client}
}

// CopyFileClient 拷贝文件客户端
type CopyFileClient struct {
	c common.Client
}

// Client 返回使用的客户端
func (c *CopyFileClient) Client() common.Client {
	return c.c
}

// WithClient 返回使用指定客户端的拷贝文件客户端
func (c *CopyFileClient) WithClient(client common.Client) *CopyFileClient {
	return &CopyFileClient{
		c: client,
	}
}

// Send 发送文件或目录
func (c *CopyFileClient) Send(ctx context.Context, stream *streamv1.Stream, path string) error {
	logger := logr.FromContextOrDiscard(ctx)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 与服务端建立连接
	conn, err := c.c.ConnectStream(ctx, stream.Name, common.ConnectStreamOptions{
		ConnectionName: "sender",
	})
	if err != nil {
		return fmt.Errorf("connect to server error: %w", err)
	}
	if logger.V(1).Enabled() {
		conn = streams.ConnectionWithLog{Connection: conn}
	}
	defer func() {
		_ = conn.Close(ctx)
	}()

	// 等待接收端
	logger.Info("waiting for receiver ...")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		msg, err := conn.Receive(ctx)
		if err != nil {
			return fmt.Errorf("receive message error: %w", err)
		}
		if len(msg) == len(startRecvMsg) && string(msg) == startRecvMsg {
			// 开始传输
			break
		}
	}
	logger.Info("start send")

	pipeR, pipeW := io.Pipe()
	defer func() {
		_ = pipeR.Close()
	}()

	// 打包文件
	tarW := tar.NewWriter(pipeW)
	go func() {
		if err := c.tarFiles(ctx, path, tarW); err != nil {
			logger.Error(err, "tar files error")
		}
		_ = tarW.Close()
		_ = pipeW.Close()
	}()

	// 转发到服务端
	tmp := make([]byte, maxReadSize)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, err := pipeR.Read(tmp)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if err := conn.Send(ctx, tmp[:n]); err != nil {
			return fmt.Errorf("send to server error: %w", err)
		}
	}

	// 发送结束消息
	if err := conn.Send(ctx, []byte(sendDoneMsg)); err != nil {
		return fmt.Errorf("send to server error: %w", err)
	}
	// 等待接收端结束
	logger.Info("send completed, waiting for receiver...")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		msg, err := conn.Receive(ctx)
		if err != nil {
			if errors.Is(err, streams.ErrConnectionClosed) {
				// 接收结束
				break
			}
			return fmt.Errorf("receive message error: %w", err)
		}
		if len(msg) == len(recvDoneMsg) && string(msg) == recvDoneMsg {
			// 接收结束
			break
		}
	}
	logger.Info("receive completed")

	return nil
}

// Receive 接收文件或目录
func (c *CopyFileClient) Receive(ctx context.Context, stream *streamv1.Stream, path string) (string, error) {
	logger := logr.FromContextOrDiscard(ctx)

	trimName := false
	stat, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
		// 文件不存在
		trimName = true
	} else if !stat.IsDir() {
		// 是个文件
		return "", fmt.Errorf("path %s is a file", path)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 与服务端建立连接
	conn, err := c.c.ConnectStream(ctx, stream.Name, common.ConnectStreamOptions{
		ConnectionName: "receiver",
	})
	if err != nil {
		return "", fmt.Errorf("connect to server error: %w", err)
	}
	if logger.V(1).Enabled() {
		conn = streams.ConnectionWithLog{Connection: conn}
	}
	defer func() {
		_ = conn.Close(ctx)
	}()

	// 发送开始传输指令
	if err := conn.Send(ctx, []byte(startRecvMsg)); err != nil {
		return "", fmt.Errorf("send start message error: %w", err)
	}

	pipeR, pipeW := io.Pipe()
	defer func() {
		_ = pipeW.Close()
	}()

	// 解包文件
	tarR := tar.NewReader(pipeR)
	untarDone := make(chan struct{})
	var untarErr error
	untarTarget := ""
	go func() {
		untarTarget, untarErr = c.untarFiles(ctx, path, tarR, trimName)
		_ = pipeR.Close()
		close(untarDone)
	}()

	// 从服务端读
	for {
		select {
		case <-ctx.Done():
			return untarTarget, ctx.Err()
		default:
		}

		data, err := conn.Receive(ctx)
		if err != nil {
			return untarTarget, fmt.Errorf("receive from server error: %w", err)
		}

		if len(data) == len(sendDoneMsg) && string(data) == sendDoneMsg {
			// 读完了
			_ = pipeW.Close()
			break
		}
		if _, err := pipeW.Write(data); err != nil {
			return untarTarget, err
		}
	}

	if err := conn.Send(ctx, []byte(recvDoneMsg)); err != nil {
		logger.Error(err, "send to server error")
	}

	// 等待解包完成
	<-untarDone

	return untarTarget, untarErr
}

// tarFiles 打包文件
func (c *CopyFileClient) tarFiles(ctx context.Context, root string, w *tar.Writer) error {
	logger := logr.FromContextOrDiscard(ctx)

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("get abs path for %q err: %v", root, err)
	}
	name := filepath.Base(absRoot)
	if name == "/" {
		name = "rootfs"
	}
	return filepath.Walk(absRoot, func(path string, info fs.FileInfo, err error) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		nameInTar := filepath.Join(name, strings.TrimPrefix(path, absRoot))

		var header *tar.Header
		mode := info.Mode()
		switch {
		case mode.IsDir():
			header = &tar.Header{
				Typeflag: tar.TypeDir,
				Name:     nameInTar,
				Mode:     int64(mode.Perm()),
			}
			logger.V(1).Info(fmt.Sprintf("dir  %s (%s)", nameInTar, mode.Perm().String()))
		case mode&fs.ModeSymlink != 0:
			link, err := os.Readlink(path)
			if err != nil {
				return fmt.Errorf("read symlink %q err: %v", path, err)
			}
			header = &tar.Header{
				Typeflag: tar.TypeSymlink,
				Name:     nameInTar,
				Linkname: link,
			}
			logger.V(1).Info(fmt.Sprintf("link %s -> %s", nameInTar, link))
		case mode.IsRegular():
			header = &tar.Header{
				Typeflag: tar.TypeReg,
				Name:     nameInTar,
				Size:     info.Size(),
				Mode:     int64(mode.Perm()),
			}
			logger.V(1).Info(fmt.Sprintf("file %s (%s)", nameInTar, mode.Perm().String()))
		default:
			logger.V(1).Info(fmt.Sprintf("skip %s", nameInTar))
			return nil
		}

		// 写文件头
		if err := w.WriteHeader(header); err != nil {
			return fmt.Errorf("write header %q to tar error: %w", nameInTar, err)
		}

		if !mode.IsRegular() {
			return nil
		}

		// 写文件
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open file %q error: %w", path, err)
		}
		defer func() {
			_ = f.Close()
		}()
		if _, err := io.Copy(w, f); err != nil {
			return fmt.Errorf("write file %q to %q in tar error: %w", path, nameInTar, err)
		}
		return nil
	})
}

// untarFiles 解包文件
func (c *CopyFileClient) untarFiles(ctx context.Context, root string, r *tar.Reader, trimName bool) (string, error) {
	logger := logr.FromContextOrDiscard(ctx)

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("get abs path for %q err: %v", root, err)
	}

	first := true
	name := ""
	target := ""
	for {
		select {
		case <-ctx.Done():
			return target, ctx.Err()
		default:
		}

		hdr, err := r.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return target, err
		}
		if first {
			name = hdr.Name
		}
		var path string
		if trimName {
			path = filepath.Join(absRoot, strings.TrimPrefix(hdr.Name, name))
		} else {
			path = filepath.Join(absRoot, hdr.Name)
		}
		if first {
			first = false
			target = path
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			logger.V(1).Info(fmt.Sprintf("dir  %s (%s)", path, hdr.FileInfo().Mode()))
			if err := os.MkdirAll(path, hdr.FileInfo().Mode()); err != nil {
				return target, fmt.Errorf("mkdir %q error: %w", path, err)
			}
		case tar.TypeSymlink:
			logger.V(1).Info(fmt.Sprintf("link %s -> %s", path, hdr.Linkname))
			if err := os.Symlink(path, hdr.Linkname); err != nil {
				return target, fmt.Errorf("symlink %q error: %w", path, err)
			}
		case tar.TypeReg:
			logger.V(1).Info(fmt.Sprintf("file %s (%s)", path, hdr.FileInfo().Mode()))
			f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, hdr.FileInfo().Mode())
			if err != nil {
				_ = f.Close()
				return target, fmt.Errorf("open file %q error: %w", path, err)
			}
			if _, err := io.Copy(f, r); err != nil {
				_ = f.Close()
				return target, fmt.Errorf("write file %q error: %w", path, err)
			}
			_ = f.Close()
		default:
			logger.V(1).Info(fmt.Sprintf("skip %s", hdr.Name))
			continue
		}
	}

	return target, nil
}
