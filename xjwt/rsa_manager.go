package xjwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// RSAManagerBuilder RSA jwt 管理器 builder。
// 注意默认 token 过期时间为 24 小时。
type RSAManagerBuilder[T any] struct {
	config ClaimsConfig

	encryptKey string
	decryptKey string
}

func NewRSAManagerBuilder[T any](encryptKey, decryptKey string) *RSAManagerBuilder[T] {
	return &RSAManagerBuilder[T]{
		config:     NewClaimsConfig(), // 默认 24 小时过期
		encryptKey: encryptKey,
		decryptKey: decryptKey,
	}
}

func NewRSAVerifierBuilder[T any](decryptKey string) *RSAManagerBuilder[T] {
	return &RSAManagerBuilder[T]{
		config:     NewClaimsConfig(),
		decryptKey: decryptKey,
	}
}

func (b *RSAManagerBuilder[T]) Build() (*RSAManager[T], error) {
	var (
		priKey *rsa.PrivateKey
		err    error
	)
	if b.encryptKey != "" {
		priKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(b.encryptKey))
		if err != nil {
			return nil, fmt.Errorf("[jit] failed to parse private key: %w", err)
		}
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(b.decryptKey))
	if err != nil {
		return nil, fmt.Errorf("[jit] failed to parse public key: %w", err)
	}

	return &RSAManager[T]{
		config: b.config,
		priKey: priKey,
		pubKey: pubKey,
	}, nil
}

func (b *RSAManagerBuilder[T]) ClaimsConfig(config ClaimsConfig) *RSAManagerBuilder[T] {
	b.config = config
	return b
}

var _ Manager[any] = (*RSAManager[any])(nil)

type RSAManager[T any] struct {
	config ClaimsConfig

	priKey *rsa.PrivateKey // 私钥
	pubKey *rsa.PublicKey  // 公钥
}

// Sign 签发一个携带 data 的 jwt token。
func (m *RSAManager[T]) Sign(data T) (string, error) {
	if m.priKey == nil {
		return "", fmt.Errorf("[jit] private key not available")
	}
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

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, cc)
	return token.SignedString(m.priKey)
}

// Verify 校验 token 签名并解析 claims。
func (m *RSAManager[T]) Verify(token string, opts ...jwt.ParserOption) (CustomClaims[T], error) {
	jwtToken, err := jwt.ParseWithClaims(
		token,
		&CustomClaims[T]{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("[jit] unexpected signing method: %v", token.Header["alg"])
			}
			return m.pubKey, nil
		},
		opts...,
	)
	if err != nil {
		return CustomClaims[T]{}, fmt.Errorf("[jit] failed to verify jwt token: %w", err)
	}
	if !jwtToken.Valid {
		return CustomClaims[T]{}, errors.New("[jit] invalid jwt token")
	}
	cc, _ := jwtToken.Claims.(*CustomClaims[T])
	return *cc, nil
}

// Encrypt 签发一个携带 data 的 jwt token。
//
// Deprecated: jwt 是签名而非加密，
// 请使用 Sign，将在未来版本移除。
func (m *RSAManager[T]) Encrypt(data T) (string, error) {
	return m.Sign(data)
}

// Decrypt 校验 token 签名并解析 claims。
//
// Deprecated: jwt 是签名而非加密，
// 请使用 Verify，将在未来版本移除。
func (m *RSAManager[T]) Decrypt(token string, opts ...jwt.ParserOption) (CustomClaims[T], error) {
	return m.Verify(token, opts...)
}
