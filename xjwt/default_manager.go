package xjwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// DefaultManagerBuilder 默认 jwt 管理器 builder。
// 注意默认 token 过期时间为 24 小时。
type DefaultManagerBuilder[T any] struct {
	config ClaimsConfig

	signingMethod jwt.SigningMethod
	encryptKey    string
	decryptKey    string
}

func NewDefaultManagerBuilder[T any](encryptKey, decryptKey string) *DefaultManagerBuilder[T] {
	return &DefaultManagerBuilder[T]{
		config: NewClaimsConfig(),

		signingMethod: jwt.SigningMethodHS256,
		encryptKey:    encryptKey,
		decryptKey:    decryptKey,
	}
}

// NewDefaultVerifierBuilder 创建仅验证场景的默认 jwt 管理器 builder，
// 不指定具体算法，Verify 接受 HMAC 系 ( HS256/HS384/HS512 ) 签发的 token。
func NewDefaultVerifierBuilder[T any](decryptKey string) *DefaultManagerBuilder[T] {
	return &DefaultManagerBuilder[T]{
		config:     NewClaimsConfig(),
		decryptKey: decryptKey,
	}
}

func (b *DefaultManagerBuilder[T]) Build() *DefaultManager[T] {
	return &DefaultManager[T]{
		config:        b.config,
		signingMethod: b.signingMethod,
		encryptKey:    b.encryptKey,
		decryptKey:    b.decryptKey,
	}
}

func (b *DefaultManagerBuilder[T]) ClaimsConfig(config ClaimsConfig) *DefaultManagerBuilder[T] {
	b.config = config
	return b
}

func (b *DefaultManagerBuilder[T]) SigningMethod(signingMethod jwt.SigningMethod) *DefaultManagerBuilder[T] {
	b.signingMethod = signingMethod
	return b
}

var _ Manager[any] = (*DefaultManager[any])(nil)

type DefaultManager[T any] struct {
	config ClaimsConfig

	signingMethod jwt.SigningMethod
	encryptKey    string
	decryptKey    string
}

// Sign 签发一个携带 data 的 jwt token。
func (m *DefaultManager[T]) Sign(data T) (string, error) {
	token := jwt.NewWithClaims(m.signingMethod, signClaims(m.config, data, time.Now()))
	if err := applyTokenHeader(token, m.config.Header); err != nil {
		return "", err
	}

	return token.SignedString([]byte(m.encryptKey))
}

// Verify 校验 token 签名并解析 claims。
// 配置了 Audience 时按该受众强制校验。
func (m *DefaultManager[T]) Verify(token string, opts ...jwt.ParserOption) (CustomClaims[T], error) {
	jwtToken, err := jwt.ParseWithClaims(
		token,
		newClaims[T](m.config),
		func(token *jwt.Token) (any, error) {
			if err := m.checkSigningMethod(token); err != nil {
				return nil, err
			}
			return []byte(m.decryptKey), nil
		},
		withConfiguredAudience(m.config, opts)...,
	)
	if err != nil {
		return CustomClaims[T]{}, fmt.Errorf("[jit] failed to verify jwt token: %w", err)
	}
	if !jwtToken.Valid {
		return CustomClaims[T]{}, errors.New("[jit] invalid jwt token")
	}

	return toCustomClaims[T](jwtToken.Claims)
}

// checkSigningMethod 拒绝非预期算法签发的 token，防止算法混淆攻击。
// 密钥为 HMAC 密钥，因此一律要求 HMAC 系算法；
// 配置了具体算法时 ( 签发场景 ) 进一步要求精确匹配。
func (m *DefaultManager[T]) checkSigningMethod(token *jwt.Token) error {
	unexpectedErr := fmt.Errorf("[jit] unexpected signing method: %v", token.Header["alg"])

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return unexpectedErr
	}

	if m.signingMethod != nil && token.Method.Alg() != m.signingMethod.Alg() {
		return unexpectedErr
	}

	return nil
}

// Encrypt 签发一个携带 data 的 jwt token。
//
// Deprecated: jwt 是签名而非加密，
// 请使用 Sign，将在未来版本移除。
func (m *DefaultManager[T]) Encrypt(data T) (string, error) {
	return m.Sign(data)
}

// Decrypt 校验 token 签名并解析 claims。
//
// Deprecated: jwt 是签名而非加密，
// 请使用 Verify，将在未来版本移除。
func (m *DefaultManager[T]) Decrypt(token string, opts ...jwt.ParserOption) (CustomClaims[T], error) {
	return m.Verify(token, opts...)
}
