package xjwt

import (
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
	const expiration = 24 * time.Hour
	return &DefaultManagerBuilder[T]{
		config: NewClaimsConfig(WithExpiration(expiration)),

		signingMethod: jwt.SigningMethodHS256,
		encryptKey:    encryptKey,
		decryptKey:    decryptKey,
	}
}

func NewDefaultVerifierBuilder[T any](decryptKey string) *DefaultManagerBuilder[T] {
	const expiration = 24 * time.Hour
	return &DefaultManagerBuilder[T]{
		config:     NewClaimsConfig(WithExpiration(expiration)),
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

func (m *DefaultManager[T]) Encrypt(data T) (string, error) {
	now := time.Now()
	cc := &CustomClaims[T]{
		Data: data,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.config.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.config.Expiration)),
			ID:        m.config.JtiGenerator(),
		},
	}

	token := jwt.NewWithClaims(m.signingMethod, cc)
	return token.SignedString([]byte(m.encryptKey))
}

func (m *DefaultManager[T]) Decrypt(token string, opts ...jwt.ParserOption) (CustomClaims[T], error) {
	jwtToken, err := jwt.ParseWithClaims(
		token,
		&CustomClaims[T]{},
		func(_ *jwt.Token) (any, error) {
			return []byte(m.decryptKey), nil
		},
		opts...,
	)
	if err != nil || !jwtToken.Valid {
		return CustomClaims[T]{}, fmt.Errorf("[jit] failed to verify jwt token: %w", err)
	}
	cc, _ := jwtToken.Claims.(*CustomClaims[T])
	return *cc, nil
}
