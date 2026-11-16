package xjwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// buildRegisteredClaims 按配置构造 RFC 7519 注册 claims。
func (cfg ClaimsConfig) buildRegisteredClaims(now time.Time) jwt.RegisteredClaims {
	claims := jwt.RegisteredClaims{
		Issuer:    cfg.Issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(cfg.Expiration)),
		ID:        cfg.JtiGenerator(),
	}

	if len(cfg.Audience) > 0 {
		claims.Audience = jwt.ClaimStrings(cfg.Audience)
	}

	return claims
}

// newClaims 构造与配置布局匹配的空 claims 对象，用于解析。
func newClaims[T any](cfg ClaimsConfig) jwt.Claims {
	if cfg.FlatClaims {
		return &FlatClaims[T]{}
	}
	return &CustomClaims[T]{}
}

// signClaims 构造携带 data 的 claims 对象，用于签发。
func signClaims[T any](cfg ClaimsConfig, data T, now time.Time) jwt.Claims {
	registered := cfg.buildRegisteredClaims(now)
	if cfg.FlatClaims {
		return &FlatClaims[T]{RegisteredClaims: registered, Data: data}
	}
	return &CustomClaims[T]{RegisteredClaims: registered, Data: data}
}

// toCustomClaims 把内部 claims 对象统一转换回 CustomClaims 返回给调用方。
func toCustomClaims[T any](claims jwt.Claims) (CustomClaims[T], error) {
	switch c := claims.(type) {
	case *CustomClaims[T]:
		return *c, nil
	case *FlatClaims[T]:
		return CustomClaims[T]{RegisteredClaims: c.RegisteredClaims, Data: c.Data}, nil
	default:
		return CustomClaims[T]{}, errors.New("[jit] unexpected claims type")
	}
}

// applyTokenHeader 把配置的自定义 header ( 如 kid ) 写入 token。
// alg 始终由签名方法决定，不允许覆盖。
func applyTokenHeader(token *jwt.Token, header map[string]any) error {
	for key := range header {
		if key == "alg" {
			return errors.New("[jit] header key \"alg\" is reserved and cannot be overridden")
		}
	}

	for key, val := range header {
		token.Header[key] = val
	}

	return nil
}

// withConfiguredAudience 在配置了 Audience 时追加等价的校验选项，
// 与网关侧 jwt.WithAudience 行为一致；
// 调用方显式传入的 opts 后写覆盖，优先生效。
func withConfiguredAudience(cfg ClaimsConfig, opts []jwt.ParserOption) []jwt.ParserOption {
	if len(cfg.Audience) == 0 {
		return opts
	}

	return append([]jwt.ParserOption{jwt.WithAudience(cfg.Audience...)}, opts...)
}
