package options

import "github.com/spf13/pflag"

// NewDefaultServeOptions 创建默认 serve 子命令选项
func NewDefaultServeOptions() ServeOptions {
	return ServeOptions{
		ListenAddr: ":9443",
		JWTIssuer:  "scaf-server",
		JWTKey:     nil,
	}
}

// ServeOptions serve 子命令选项
type ServeOptions struct {
	ListenAddr string `json:"listenAddr,omitempty" yaml:"listenAddr,omitempty"`
	JWTIssuer  string `json:"jwtIssuer,omitempty" yaml:"jwtIssuer,omitempty"`
	JWTKey     []byte `json:"jwtKey,omitempty" yaml:"jwtKey,omitempty"`
}

// AddPFlags 绑定选项到参数
func (opts *ServeOptions) AddPFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&opts.ListenAddr, "listen", "l", opts.ListenAddr, "Listen address")
	fs.StringVar(&opts.JWTIssuer, "jwt-issuer", opts.JWTIssuer, "JWT issuer name")
	fs.BytesBase64Var(&opts.JWTKey, "jwt-key", opts.JWTKey, "JWT signing key")
}
