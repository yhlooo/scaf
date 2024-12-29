[简体中文](README_CN.md) | **[English](README.md)**

**Most of this article is translated from [README_CN.md](README_CN.md) by [ChatGPT](https://chatgpt.com/). :)**

---

![GitHub License](https://img.shields.io/github/license/yhlooo/scaf)
[![GitHub Release](https://img.shields.io/github/v/release/yhlooo/scaf)](https://github.com/yhlooo/scaf/releases/latest)
[![release](https://github.com/yhlooo/scaf/actions/workflows/release.yaml/badge.svg)](https://github.com/yhlooo/scaf/actions/workflows/release.yaml)

# Scaf

Establishes point-to-point streams by pairing and relaying with reverse connections, for remote shell, file transfer, etc.

In Scaf, a session involves three participants: the server, client A, and client B. Client A first connects to the server to create a "stream". Then, client B connects to the server to join the same "stream". Once both connections are established, the server forwards data read from client A to client B and data read from client B to client A, thus enabling point-to-point stream transmission between client A and client B.

![scaf.drawio.svg](docs/images/scaf.drawio.svg)

Scaf currently provides client implementations for remote command execution and file transfer, but its capabilities extend beyond that. By using the Scaf API, you can interact directly with the Scaf server, allowing you to transmit any data within a stream established by the server, enabling more point-to-point streaming applications.

## Installation

### Binaries

Download the executable binary from the [Releases](https://github.com/yhlooo/scaf/releases) page, extract it, and place the `scaf` file into any `$PATH` directory.

### From Sources

Requires Go 1.22. Execute the following command to download the source code and build it:

```bash
go install github.com/yhlooo/scaf/cmd/scaf@latest
```

The built binary will be located in `${GOPATH}/bin` by default. Make sure this directory is included in your `$PATH`.

## Usage

### Start the Server

```bash
scaf serve
```

By default, Scaf listens on the address `:9443`. You can specify a different address using the `-l` flag.

On the listening address, the Scaf server supports both gRPC and HTTP protocols. On the Scaf client, use the `-s` flag to specify the server address. Use `grpc://<host>:<port>` for gRPC or `http://<host>:<port>` for HTTP.

### Remote Command Execution

#### Initiated by the Monitor

The monitor creates a stream and starts the command execution session:

```bash
scaf exec-remote -it -s <SERVER_URL> -- <COMMAND> [ARGS...]
```

Once the monitor is started, it will output the stream name `<STREAM_NAME>` and its connection token `<TOKEN>`. The executor uses this information to connect and begin command execution:

```bash
scaf exec -s <SERVER_URL> --stream <STREAM_NAME> --token <TOKEN>
```

Once the executor starts, the input and output of the command will be forwarded to the monitor.

#### Initiated by the Executor

The executor creates a stream and starts the command execution session:

```bash
scaf exec -it -s <SERVER_URL> -- <COMMAND> [ARGS...]
```

Once the executor starts, it will output the stream name `<STREAM_NAME>` and its connection token `<TOKEN>`. The monitor uses this information to connect and receive the input/output:

```bash
scaf attach -s <SERVER_URL> --stream <STREAM_NAME> --token <TOKEN>
```

### File Transfer

The sender creates a stream and starts the file sending session:

```bash
scaf send-file -s <SERVER_URL> <PATH>
```

`<PATH>` refers to the file or directory being sent.

Once the sender starts, it will output the stream name `<STREAM_NAME>` and its connection token `<TOKEN>`. The receiver uses this information to connect and receive the file:

```bash
scaf receive-file -s <SERVER_URL> --stream <STREAM_NAME> --token <TOKEN> [PATH]
```

`[PATH]` is an optional path to save the received file. If not specified, the current working directory will be used.

Once the receiver connects, the file transfer will begin.
