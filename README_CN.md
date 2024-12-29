**[简体中文](README_CN.md)** | [English](README.md)

---

![GitHub License](https://img.shields.io/github/license/yhlooo/scaf)
[![GitHub Release](https://img.shields.io/github/v/release/yhlooo/scaf)](https://github.com/yhlooo/scaf/releases/latest)
[![release](https://github.com/yhlooo/scaf/actions/workflows/release.yaml/badge.svg)](https://github.com/yhlooo/scaf/actions/workflows/release.yaml)

# Scaf

Scaf 通过配对和中继反向连接建立点到点的流传输，用于远程命令执行、文件传输等。

在 Scaf 中一个会话由服务端、客户端 A 和客户端 B 三方参与。客户端 A 首先连接到服务端创建“流”，随后客户端 B 连接到服务端接入相同的“流”，两个连接建立后服务端将从客户端 A 读取的数据转发到 B 并将从 B 读取的数据转发到 A 。从而实现客户端 A 和 B 之间的点到点的流传输。

Scaf 目前提供用于远程命令执行和文件传输的客户端实现，但 Scaf 的功能不止于此，使用 Scaf API 可直接与 Scaf 服务端交互，在通过 Scaf 服务端建立的流中可以传输任何数据，基于此可以实现更多其它基于点对点的流传输的应用。

## 安装

### 通过二进制安装

通过 [Releases](https://github.com/yhlooo/scaf/releases) 页面下载可执行二进制，解压并将其中 `scaf` 文件放置到任意 `$PATH` 目录下。

### 从源码编译

要求 Go 1.22 ，执行以下命令下载源码并构建：

```bash
go install github.com/yhlooo/scaf/cmd/scaf@latest
```

构建的二进制默认将在 `${GOPATH}/bin` 目录下，需要确保该目录包含在 `$PATH` 中。

## 使用

### 启动服务

```bash
scaf serve
```

默认参数下 Scaf 监听地址为 `:9443` ，通过 `-l` 参数可指定其他地址。

在监听地址上 Scaf 服务端同时支持 gRPC 和 HTTP 协议。在 Scaf 客户端通过 `-s` 参数指定服务端地址，使用 `grpc://<host>:<port>` 指定以 gRPC 协议访问，使用 `http://<host>:<port>` 指定使用 HTTP 协议访问。  

### 远程执行命令

#### 由监视端发起

监视端创建流，开启命令执行会话：

```bash
scaf exec-remote -it -s <SERVER_URL> -- <COMMAND> [ARGS...]
```

监视端开启后会输出流名 `<STREAM_NAME>` 和其连接凭证 `<TOKEN>` ，在执行端使用该信息连接并开始命令执行：

```bash
scaf exec -s <SERVER_URL> --stream <STREAM_NAME> --token <TOKEN>
```

执行端开始执行后，命令执行的输入输出会转发到监视端。

#### 由执行端发起

执行端创建流，开启命令执行会话：

```bash
scaf exec -it -s <SERVER_URL> -- <COMMAND> [ARGS...]
```

执行端开始执行后会输出流名 <STREAM_NAME> 和其连接凭证 <TOKEN> ，在监视端使用该信息连接并获取输入输出：

```bash
scaf attach -s <SERVER_URL> --stream <STREAM_NAME> --token <TOKEN>
```

### 传输文件

在发送端创建流，开启文件发送会话：

```bash
scaf send-file -s <SERVER_URL> <PATH>
```

`<PATH>` 是发送的文件或目录路径。

发送端开启后会输出流名 `<STREAM_NAME>` 和其连接凭证 `<TOKEN>` ，在接收端使用该信息连接并获取文件：

```bash
scaf receive-file -s <SERVER_URL> --stream <STREAM_NAME> --token <TOKEN> [PATH]
```

`[PATH]` 是可选的接收文件的路径，未指定时使用当前工作目录。

接收端连接后，文件会开始传输。
