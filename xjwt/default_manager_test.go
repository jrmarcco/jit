package xjwt

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type defaultUser struct {
	ID uint64
}

func TestDefaultManager(t *testing.T) {
	t.Parallel()

	key := "test-key"
	tcs := []struct {
		name          string
		user          defaultUser
		manager       Manager[defaultUser]
		wantIssuer    string
		wantSignErr   error
		wantVerifyErr error
	}{
		{
			name: "basic",
			user: defaultUser{ID: 1},
			manager: func() *DefaultManager[defaultUser] {
				return NewDefaultManagerBuilder[defaultUser](key, key).Build()
			}(),
			wantIssuer:    defaultIssuer,
			wantSignErr:   nil,
			wantVerifyErr: nil,
		}, {
			name: "expired",
			user: defaultUser{ID: 1},
			manager: func() *DefaultManager[defaultUser] {
				return NewDefaultManagerBuilder[defaultUser](key, key).
					ClaimsConfig(NewClaimsConfig(WithExpiration(time.Millisecond))).
					Build()
			}(),
			wantIssuer:    defaultIssuer,
			wantSignErr:   nil,
			wantVerifyErr: jwt.ErrTokenExpired,
		}, {
			name: "with issuer",
			user: defaultUser{ID: 1},
			manager: func() *DefaultManager[defaultUser] {
				cfg := NewClaimsConfig(WithExpiration(time.Minute), WithIssuer("test-issuer"))
				return NewDefaultManagerBuilder[defaultUser](key, key).
					ClaimsConfig(cfg).
					Build()
			}(),
			wantIssuer:    "test-issuer",
			wantSignErr:   nil,
			wantVerifyErr: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			token, err := tc.manager.Sign(tc.user)
			assert.Equal(t, tc.wantSignErr, err)
			if err != nil {
				return
			}

			time.Sleep(time.Millisecond)
			decrypted, err := tc.manager.Verify(token)
			assert.Truef(t, errors.Is(err, tc.wantVerifyErr), "want: %v, got: %v", tc.wantVerifyErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.user, decrypted.Data)
			assert.Equal(t, tc.wantIssuer, decrypted.Issuer)
		})
	}
}

func TestDefaultManager_InvalidToken(t *testing.T) {
	t.Parallel()

	key := "test-key"
	manager := NewDefaultManagerBuilder[defaultUser](key, key).Build()

	var err error
	_, err = manager.Sign(defaultUser{ID: 1})
	require.NoError(t, err)

	_, err = manager.Verify("invalid token")
	wantErr := jwt.ErrTokenMalformed
	assert.Truef(t, errors.Is(err, wantErr), "want: %v, got: %v", wantErr, err)
}

// TestDefaultManager_DeprecatedAliases 验证废弃的
// Encrypt/Decrypt 方法仍委托给 Sign/Verify。
func TestDefaultManager_DeprecatedAliases(t *testing.T) {
	t.Parallel()

	manager := NewDefaultManagerBuilder[defaultUser]("alias-test-key", "alias-test-key").Build()

	token, err := manager.Encrypt(defaultUser{ID: 1})
	require.NoError(t, err)

	claims, err := manager.Decrypt(token)
	require.NoError(t, err)
	assert.Equal(t, defaultUser{ID: 1}, claims.Data)
}

// TestDefaultManager_AlgMismatch 验证以非预期算法签发的 token 会被拒绝，
// 即使签名密钥相同 ( 算法混淆攻击 )。
func TestDefaultManager_AlgMismatch(t *testing.T) {
	t.Parallel()

	// 使用独立 key 常量，避免推高共享 "test-key" 字符串的
	// goconst 计数。
	const algTestKey = "alg-mismatch-test-key"

	// 以 HS512 配置的 manager。
	manager := NewDefaultManagerBuilder[defaultUser](algTestKey, algTestKey).
		SigningMethod(jwt.SigningMethodHS512).
		Build()

	// 以相同 key 用 HS256 签发的 token。
	claims := &CustomClaims[defaultUser]{
		Data: defaultUser{ID: 1},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    defaultIssuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(algTestKey))
	require.NoError(t, err)

	_, err = manager.Verify(token)
	assert.Error(t, err, "token signed with a mismatched algorithm must be rejected")
}
