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
		name           string
		user           defaultUser
		manager        Manager[defaultUser]
		wantIssuer     string
		wantEncryptErr error
		wantDecryptErr error
	}{
		{
			name: "basic",
			user: defaultUser{ID: 1},
			manager: func() *DefaultManager[defaultUser] {
				return NewDefaultManagerBuilder[defaultUser](key, key).Build()
			}(),
			wantIssuer:     "jit",
			wantEncryptErr: nil,
			wantDecryptErr: nil,
		}, {
			name: "expired",
			user: defaultUser{ID: 1},
			manager: func() *DefaultManager[defaultUser] {
				return NewDefaultManagerBuilder[defaultUser](key, key).
					ClaimsConfig(NewClaimsConfig(WithExpiration(time.Millisecond))).
					Build()
			}(),
			wantIssuer:     "jit",
			wantEncryptErr: nil,
			wantDecryptErr: jwt.ErrTokenExpired,
		}, {
			name: "with issuer",
			user: defaultUser{ID: 1},
			manager: func() *DefaultManager[defaultUser] {
				cfg := NewClaimsConfig(WithExpiration(time.Minute), WithIssuer("test-issuer"))
				return NewDefaultManagerBuilder[defaultUser](key, key).
					ClaimsConfig(cfg).
					Build()
			}(),
			wantIssuer:     "test-issuer",
			wantEncryptErr: nil,
			wantDecryptErr: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			token, err := tc.manager.Encrypt(tc.user)
			assert.Equal(t, tc.wantEncryptErr, err)
			if err != nil {
				return
			}

			time.Sleep(time.Millisecond)
			decrypted, err := tc.manager.Decrypt(token)
			assert.Truef(t, errors.Is(err, tc.wantDecryptErr), "want: %v, got: %v", tc.wantDecryptErr, err)
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
	_, err = manager.Encrypt(defaultUser{ID: 1})
	require.NoError(t, err)

	_, err = manager.Decrypt("invalid token")
	wantErr := jwt.ErrTokenMalformed
	assert.Truef(t, errors.Is(err, wantErr), "want: %v, got: %v", wantErr, err)
}
