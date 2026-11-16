// Package xjwt 提供基于 github.com/golang-jwt/jwt/v5 的 jwt 管理器，
// 支持 HMAC、RSA 与 Ed25519 签名算法，泛型承载自定义 claims 数据，
// 支持 aud claim、自定义 token header ( 如 kid ) 与顶层平铺 claims 布局。
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
	Issuer       string         // 签发人
	Audience     []string       // 受众 ( RFC 7519 aud claim )；非空时 Verify 会强制校验
	Expiration   time.Duration  // 有效期
	JtiGenerator func() string  // jwt id 生成方法
	FlatClaims   bool           // 业务数据平铺到 payload 顶层 ( FlatClaims 布局 )
	Header       map[string]any // 签发时写入 token header 的自定义字段 ( 如 kid )
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

// WithAudience 设置受众 ( RFC 7519 aud claim )。
// 配置后签发的 token 携带 aud，Verify 亦按该受众强制校验；
// 调用方在 Verify 显式传入的 jwt.WithAudience 优先生效。
func WithAudience(audience ...string) option.Opt[ClaimsConfig] {
	return func(cfg *ClaimsConfig) {
		cfg.Audience = audience
	}
}

func WithJtiGenerator(jtiGenerator func() string) option.Opt[ClaimsConfig] {
	return func(cfg *ClaimsConfig) {
		cfg.JtiGenerator = jtiGenerator
	}
}

// WithFlatClaims 启用顶层平铺 claims 布局 ( FlatClaims[T] )，
// 业务数据的 JSON 键直接写入 payload 顶层，而非 "Data" 键下。
func WithFlatClaims() option.Opt[ClaimsConfig] {
	return func(cfg *ClaimsConfig) {
		cfg.FlatClaims = true
	}
}

// WithHeader 设置签发时写入 token header 的自定义字段，
// 如 JWKS 密钥轮换所需的 kid。
// "alg" 由签名方法决定，不允许覆盖。
func WithHeader(header map[string]any) option.Opt[ClaimsConfig] {
	return func(cfg *ClaimsConfig) {
		cfg.Header = header
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
