package xjwt

import (
	"crypto/rsa"
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

func (b *RSAManagerBuilder[T]) ClaimsConfig(config ClaimsConfig) *RSAManagerBuilder[T] {
	b.config = config
	return b
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

func NewRSAManagerBuilder[T any](encryptKey, decryptKey string) *RSAManagerBuilder[T] {
	const expiration = 24 * time.Hour
	return &RSAManagerBuilder[T]{
		config:     NewClaimsConfig(WithExpiration(expiration)), // 默认 24 小时过期
		encryptKey: encryptKey,
		decryptKey: decryptKey,
	}
}

func NewRSAVerifierBuilder[T any](decryptKey string) *RSAManagerBuilder[T] {
	const expiration = 24 * time.Hour
	return &RSAManagerBuilder[T]{
		config:     NewClaimsConfig(WithExpiration(expiration)),
		decryptKey: decryptKey,
	}
}

var _ Manager[any] = (*RSAManager[any])(nil)

type RSAManager[T any] struct {
	config ClaimsConfig

	priKey *rsa.PrivateKey // 私钥
	pubKey *rsa.PublicKey  // 公钥
}

func (m *RSAManager[T]) Encrypt(data T) (string, error) {
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

func (m *RSAManager[T]) Decrypt(token string, opts ...jwt.ParserOption) (CustomClaims[T], error) {
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
	if err != nil || !jwtToken.Valid {
		return CustomClaims[T]{}, fmt.Errorf("[jit] failed to verify jwt token: %w", err)
	}
	cc, _ := jwtToken.Claims.(*CustomClaims[T])
	return *cc, nil
}
