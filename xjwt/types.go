// Package xjwt 提供基于 github.com/golang-jwt/jwt/v5 的 jwt 管理器，
// 支持 HMAC、RSA 与 Ed25519 签名算法，泛型承载自定义 claims 数据。
package xjwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jrmarcco/jit/xbean/option"
)

const (
	// defaultExpiration 为统一的默认 token 有效期，
	// NewClaimsConfig 与各 manager builder 共用。
	defaultExpiration = 24 * time.Hour
	// defaultIssuer 为统一的默认签发人。
	defaultIssuer = "jit"
)

// Manager jwt 管理器抽象。
// jwt 是签名机制而非加密机制，因此方法命名为 Sign / Verify。
type Manager[T any] interface {
	// Sign 签发一个携带 data 的 jwt token。
	Sign(data T) (string, error)
	// Verify 校验 token 签名与有效性，并解析出 claims。
	Verify(token string, opts ...jwt.ParserOption) (CustomClaims[T], error)
}

type CustomClaims[T any] struct {
	jwt.RegisteredClaims
	Data T
}

// ClaimsConfig jwt claims 扩展配置项
type ClaimsConfig struct {
	Issuer       string        // 签发人
	Expiration   time.Duration // 有效期
	JtiGenerator func() string // jwt id 生成方法
}

func WithExpiration(expiration time.Duration) option.Opt[ClaimsConfig] {
	return func(cfg *ClaimsConfig) {
		cfg.Expiration = expiration
	}
}

func WithIssuer(issuer string) option.Opt[ClaimsConfig] {
	return func(cfg *ClaimsConfig) {
		cfg.Issuer = issuer
	}
}

func WithJtiGenerator(jtiGenerator func() string) option.Opt[ClaimsConfig] {
	return func(cfg *ClaimsConfig) {
		cfg.JtiGenerator = jtiGenerator
	}
}

// NewClaimsConfig 创建自定义 claims 配置信息。
// 其他参数可以通过 option.Opt 自定义。
func NewClaimsConfig(opts ...option.Opt[ClaimsConfig]) ClaimsConfig {
	cfg := ClaimsConfig{
		Issuer:       defaultIssuer, // 默认签发人
		Expiration:   defaultExpiration,
		JtiGenerator: func() string { return "" },
	}

	option.Apply(&cfg, opts...)

	return cfg
}
