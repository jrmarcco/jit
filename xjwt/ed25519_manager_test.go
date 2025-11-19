package xjwt

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	priPem = `-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEIHXEAUN6Lp8Hdq8P0Mcv9mjIG1sgPWBf1Mh+OKP5HXvC
-----END PRIVATE KEY-----`
	pubPem = `-----BEGIN PUBLIC KEY-----
MCowBQYDK2VwAyEAxCSxEyY/+A7T7EtXF7AHw4Zfklh/QdjG8fxfRFYZgY8=
-----END PUBLIC KEY-----`
)

type ed25519User struct {
	ID uint64
}

func TestEd25519Manager(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name           string
		user           ed25519User
		manager        Manager[ed25519User]
		wantIssuer     string
		wantEncryptErr error
		wantDecryptErr error
	}{
		{
			name: "basic",
			user: ed25519User{ID: 1},
			manager: func() *Ed25519Manager[ed25519User] {
				manager, err := NewEd25519ManagerBuilder[ed25519User](priPem, pubPem).Build()
				require.NoError(t, err)
				return manager
			}(),
			wantIssuer:     "jit",
			wantEncryptErr: nil,
			wantDecryptErr: nil,
		}, {
			name: "expired",
			user: ed25519User{ID: 1},
			manager: func() *Ed25519Manager[ed25519User] {
				manager, err := NewEd25519ManagerBuilder[ed25519User](priPem, pubPem).
					ClaimsConfig(NewClaimsConfig(WithExpiration(time.Millisecond))).
					Build()
				require.NoError(t, err)
				return manager
			}(),
			wantIssuer:     "jit",
			wantEncryptErr: nil,
			wantDecryptErr: jwt.ErrTokenExpired,
		}, {
			name: "with issuer",
			user: ed25519User{ID: 1},
			manager: func() *Ed25519Manager[ed25519User] {
				cfg := NewClaimsConfig(WithExpiration(time.Minute), WithIssuer("test-issuer"))
				manager, err := NewEd25519ManagerBuilder[ed25519User](priPem, pubPem).
					ClaimsConfig(cfg).
					Build()
				require.NoError(t, err)
				return manager
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

func TestEd25519Manager_InvalidToken(t *testing.T) {
	t.Parallel()

	manager, err := NewEd25519ManagerBuilder[ed25519User](priPem, pubPem).Build()
	require.NoError(t, err)

	_, err = manager.Encrypt(ed25519User{ID: 1})
	require.NoError(t, err)

	_, err = manager.Decrypt("invalid token")
	wantErr := jwt.ErrTokenMalformed
	assert.Truef(t, errors.Is(err, wantErr), "want: %v, got: %v", wantErr, err)
}
