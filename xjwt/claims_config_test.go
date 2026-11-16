package xjwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ClaimsConfig 扩展能力 ( aud / kid header / 顶层平铺 claims )
// 场景测试共用的常量。
const (
	scenarioAudience = "synp-gateway"
	scenarioKid      = "rsa-key-2026-01"
)

// 仅用于测试的固定 RSA 密钥对 ( 2048 bit )，不用于任何真实环境。
var (
	rsaPriPem = `-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQCyZghSDqpOAozF
AAG7UwJ9EiBa99F+p7QSSFrlE9v9tndMhikRBWvAkAe6deKcsSxc7vG/J5gXtthE
EfbYCnqqw2v+vGVl/dPs4R+ma8NsMG4VZTaPqZq5oEZipZZhtwRA6DW+wVn8+f4R
s33TrCK9essJ7yTKLRRlKwbK/oU6jlUI8AcpJdgOg4xYp714Cg/OMUpx3FxizeNe
r+ONg5GrkybutXaXURimebRd80suVvA0FF+HfTv1tUrD1CTFhDcnZx7evbPW0KKG
p9jiR3MxQLkcH7ZoqoOzDQc5fq860G0FwC0t1WDKI3kbtSbhnjJl62hjcRNMWMdR
YvHmW2afAgMBAAECggEAFkU72gKQABe8EFu6hkgL53vRI4oMUmNdIhXQ7bKs+TSY
xVTPi5sW+H2tI/nTcFGDLZbRVwD6W6YtGVJPc74wfLYopS70+ZErrkaW94+Rtc/V
7kB/BIWHwI825NmxtuvB/C2rh+9bQdlzfizfaFf0051LUNIp0D4QbbtmqFtdsRFz
Bxm3ZMiadcrogP/mhFU9riuJUENvgsKVdEHberr79r2gYICgPKiXf6muvM1k5sFQ
OZ2j6iAT3kt9Cry1tOwNkoDuL7cPDhHUAMPeIbJNkfNVFVaOMLGRqXZXhvMROUFU
EV6DpNx8tThLTJzCuKOpafscL8hsqBwgAreL5d/VAQKBgQD3m6kpenf2f5iG5IQn
6SGpdBMKq7JlwswF8EPWYy4/Nny0+b+iykam5PT3ftLWUPmYMhAmPeDNQvEbNHA8
rJszIvnAmu6MMf5jlIHtlEkH22sr6aQlKDaLm4AnwFmB5vOa3JmBzMoY3rej2AnA
QtavMoDK6msjFePYzBz1+/Z6JQKBgQC4ceJy7G0nXa+p9mEYmLhGNB0uNz7/cWGQ
84iYyYFWPeH83OXIia1cIvEYdHEsztFxiGzFOiIAmxbFuk7diFVphyTrQvo/2TH4
8LvQw0y191mA6QX4fjUlo9DuPmeCDYZKu8LiAurbOYiFVxwRgeAlpxelz1Mmbuwq
opNLbt3ocwKBgQCRxaOelso/ZbkvWSNwt7CoqXx4YMLjY/5nVJ8rfTKB5pselJY8
yac0pECP4KtRn1SCRJ77VJyfW1TE9mdJYTjxrxqoSM493UTt7UVc0Fr6dDQxIxkP
AxCxujeFr5lCf+xU9KCvKRil8aDdpdd/yBnKBQqL69xtSHrebLR0IQVnvQKBgQCQ
SWbgiFuLZBq7VuX2G+Jr6PC66rT5l3S/tAOzok7vbFGRQizvh+R6I0t/2fd/Gft2
WUbdsvBOkQRrgaeYXcKTXwlekXK65+nT2F7UnTMxVmJKGgB1uCVTYKdc+Kx4bKQV
anS14ZpgDM3joZVvjK4QRnqVV26wri1J3QxlAkHOSQKBgQCtXPM6A6gmPB6TK9ZC
az+MLwxZiVQGRrSvxfuHWR8NtnI9stgd6l0fAFU7SZf1zevsADqxHZKzAA918zeH
RAG1XXmrGOnWPWMkipONpk3VFo7yvuXQDREh1GaxGTVyXwN0ezHSKv1RchTqwoOm
4R61sA0G7cXJvI0PZdFhUXpm8Q==
-----END PRIVATE KEY-----`

	rsaPubPem = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsmYIUg6qTgKMxQABu1MC
fRIgWvfRfqe0Ekha5RPb/bZ3TIYpEQVrwJAHunXinLEsXO7xvyeYF7bYRBH22Ap6
qsNr/rxlZf3T7OEfpmvDbDBuFWU2j6mauaBGYqWWYbcEQOg1vsFZ/Pn+EbN906wi
vXrLCe8kyi0UZSsGyv6FOo5VCPAHKSXYDoOMWKe9eAoPzjFKcdxcYs3jXq/jjYOR
q5Mm7rV2l1EYpnm0XfNLLlbwNBRfh3079bVKw9QkxYQ3J2ce3r2z1tCihqfY4kdz
MUC5HB+2aKqDsw0HOX6vOtBtBcAtLdVgyiN5G7Um4Z4yZetoY3ETTFjHUWLx5ltm
nwIDAQAB
-----END PUBLIC KEY-----`
)

// TestRSAManager_AudienceHeaderFlatClaims 端到端验证 RS256 签名 +
// kid header + aud claim + 业务 claims 顶层平铺的组合场景
// ( 需求背景见 doc/xjwt-upgrade.md )。
func TestRSAManager_AudienceHeaderFlatClaims(t *testing.T) {
	t.Parallel()

	signer, err := NewRSAManagerBuilder[flatPayload](rsaPriPem, rsaPubPem).
		ClaimsConfig(NewClaimsConfig(
			WithAudience(scenarioAudience),
			WithExpiration(time.Hour),
			WithFlatClaims(),
			WithHeader(map[string]any{"kid": scenarioKid}),
		)).
		Build()
	require.NoError(t, err)

	want := flatPayload{BusinessID: "rsa-biz-1", UserID: 42, SessionID: "rsa-sess-1"}
	token, err := signer.Sign(want)
	require.NoError(t, err)

	// 1) header: RS256 + kid。
	unverified, _, err := jwt.NewParser().ParseUnverified(token, jwt.MapClaims{})
	require.NoError(t, err)

	assert.Equal(t, jwt.SigningMethodRS256.Alg(), unverified.Header["alg"])
	assert.Equal(t, scenarioKid, unverified.Header["kid"])

	// 2) payload: 业务 claims 平铺在顶层，无 "Data" 包装键。
	mc, ok := unverified.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, []any{scenarioAudience}, mc["aud"])
	assert.Equal(t, "rsa-biz-1", mc["business_id"])
	assert.Equal(t, float64(42), mc["user_id"])
	assert.Equal(t, "rsa-sess-1", mc["session_id"])

	_, hasDataKey := mc["Data"]
	assert.False(t, hasDataKey)

	// 3) 验证侧 round-trip: 签名、布局与受众全部通过。
	claims, err := signer.Verify(token)
	require.NoError(t, err)

	assert.Equal(t, want, claims.Data)
	assert.Equal(t, jwt.ClaimStrings{scenarioAudience}, claims.Audience)

	// 4) 受众不匹配的 manager 必须拒绝该 token。
	otherAudManager, err := NewRSAVerifierBuilder[flatPayload](rsaPubPem).
		ClaimsConfig(NewClaimsConfig(WithAudience("other-aud"), WithFlatClaims())).
		Build()
	require.NoError(t, err)

	_, err = otherAudManager.Verify(token)
	assert.Error(t, err)
}

// TestWithHeaderAlgReserved 验证 header 不允许覆盖 alg，
// alg 始终由签名方法决定。
func TestWithHeaderAlgReserved(t *testing.T) {
	t.Parallel()

	manager, err := NewRSAManagerBuilder[flatPayload](rsaPriPem, rsaPubPem).
		ClaimsConfig(NewClaimsConfig(WithHeader(map[string]any{"alg": "none"}))).
		Build()
	require.NoError(t, err)

	_, err = manager.Sign(flatPayload{})
	assert.Error(t, err)
}

// TestDefaultManager_Audience HMAC 场景下 aud 的签发与校验。
func TestDefaultManager_Audience(t *testing.T) {
	t.Parallel()

	const hmacKey = "aud-test-key"

	manager := NewDefaultManagerBuilder[flatPayload](hmacKey, hmacKey).
		ClaimsConfig(NewClaimsConfig(WithAudience(scenarioAudience))).
		Build()

	want := flatPayload{BusinessID: "h-biz-1", UserID: 1, SessionID: "h-sess-1"}
	token, err := manager.Sign(want)
	require.NoError(t, err)

	claims, err := manager.Verify(token)
	require.NoError(t, err)

	assert.Equal(t, want, claims.Data)
	assert.Equal(t, jwt.ClaimStrings{scenarioAudience}, claims.Audience)

	// 调用方显式传入的受众校验优先生效。
	_, err = manager.Verify(token, jwt.WithAudience("other-aud"))
	assert.Error(t, err)
}
