package xjwt

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Ed25519ManagerBuilder ed25519 jwt 管理器 builder。
// 注意默认 token 过期时间为 24 小时。
type Ed25519ManagerBuilder[T any] struct {
	config ClaimsConfig

	encryptKey string
	decryptKey string
}

func (b *Ed25519ManagerBuilder[T]) ClaimsConfig(config ClaimsConfig) *Ed25519ManagerBuilder[T] {
	b.config = config
	return b
}

func (b *Ed25519ManagerBuilder[T]) Build() (*Ed25519Manager[T], error) {
	var (
		priKey ed25519.PrivateKey
		err    error
	)
	if b.encryptKey != "" {
		priKey, err = loadPrivateKey(b.encryptKey)
		if err != nil {
			return nil, err
		}
	}

	pubKey, err := loadPublicKey(b.decryptKey)
	if err != nil {
		return nil, err
	}

	return &Ed25519Manager[T]{
		config: b.config,
		priKey: priKey,
		pubKey: pubKey,
	}, nil
}

// NewEd25519ManagerBuilder 创建 ed25519 管理器。
func NewEd25519ManagerBuilder[T any](encryptKey, decryptKey string) *Ed25519ManagerBuilder[T] {
	const expiration = 24 * time.Hour
	return &Ed25519ManagerBuilder[T]{
		config:     NewClaimsConfig(WithExpiration(expiration)), // 默认 24 小时过期
		encryptKey: encryptKey,
		decryptKey: decryptKey,
	}
}

// NewEd25519VerifierBuilder 用于创建只包含 public key 的 ed25519 管理器。
// 用于只需要验证 jwt 而不需要签发 jwt 的场景。
func NewEd25519VerifierBuilder[T any](decryptKey string) *Ed25519ManagerBuilder[T] {
	const expiration = 24 * time.Hour
	return &Ed25519ManagerBuilder[T]{
		config:     NewClaimsConfig(WithExpiration(expiration)),
		decryptKey: decryptKey,
	}
}

// loadPrivateKey 加载私钥。
// 注意 PEM 块本身标注的是密钥对，而不是具体的 ed25519 密钥对。
// 所有标准公钥格式都需要先由 x509 包处理进行转换后类型断言才能获得 ed25519 密钥对。
func loadPrivateKey(priPem string) (ed25519.PrivateKey, error) {
	priKeyBlock, _ := pem.Decode([]byte(priPem))
	if priKeyBlock == nil {
		return nil, fmt.Errorf("[jit] failed to decode private key PEM")
	}
	priKey, err := x509.ParsePKCS8PrivateKey(priKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("[jit] failed to parse private key: %w", err)
	}

	rtn, ok := priKey.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("[jit] private key is not ed25519")
	}
	return rtn, nil
}

// loadPublicKey 加载公钥。
func loadPublicKey(pubPem string) (ed25519.PublicKey, error) {
	pubKeyBlock, _ := pem.Decode([]byte(pubPem))
	if pubKeyBlock == nil {
		return nil, fmt.Errorf("[jit] failed to decode public key PEM")
	}
	publicKey, err := x509.ParsePKIXPublicKey(pubKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("[jit] failed to parse public key: %w", err)
	}

	rtn, ok := publicKey.(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("[jit] public key is not ed25519")
	}
	return rtn, nil
}

var _ Manager[any] = (*Ed25519Manager[any])(nil)

type Ed25519Manager[T any] struct {
	config ClaimsConfig

	priKey ed25519.PrivateKey // 私钥
	pubKey ed25519.PublicKey  // 公钥
}

func (m *Ed25519Manager[T]) Encrypt(data T) (string, error) {
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

	token := jwt.NewWithClaims(&jwt.SigningMethodEd25519{}, cc)
	return token.SignedString(m.priKey)
}

func (m *Ed25519Manager[T]) Decrypt(token string, opts ...jwt.ParserOption) (CustomClaims[T], error) {
	jwtToken, err := jwt.ParseWithClaims(
		token,
		&CustomClaims[T]{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
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
