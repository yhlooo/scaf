package options

import "github.com/spf13/pflag"

// NewDefaultServeOptions 创建默认 serve 子命令选项
func NewDefaultServeOptions() ServeOptions {
	return ServeOptions{
		HTTPAddr:  ":80",
		GRPCAddr:  ":9443",
		JWTIssuer: "scaf-server",
		JWTKey:    nil,
	}
}

// ServeOptions serve 子命令选项
type ServeOptions struct {
	HTTPAddr  string `json:"httpAddr,omitempty" yaml:"httpAddr,omitempty"`
	GRPCAddr  string `json:"grpcAddr,omitempty" yaml:"grpcAddr,omitempty"`
	JWTIssuer string `json:"jwtIssuer,omitempty" yaml:"jwtIssuer,omitempty"`
	JWTKey    []byte `json:"jwtKey,omitempty" yaml:"jwtKey,omitempty"`
}

// AddPFlags 绑定选项到参数
func (opts *ServeOptions) AddPFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&opts.HTTPAddr, "listen-http", "l", opts.HTTPAddr, "HTTP listen address")
	fs.StringVar(&opts.GRPCAddr, "listen-grpc", opts.GRPCAddr, "gRPC listen address")
	fs.StringVar(&opts.JWTIssuer, "jwt-issuer", opts.JWTIssuer, "JWT issuer name")
	fs.BytesBase64Var(&opts.JWTKey, "jwt-key", opts.JWTKey, "JWT signing key")
}
