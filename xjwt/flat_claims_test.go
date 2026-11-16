package xjwt

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// flatPayload 模拟网关契约中的顶层业务 claims ( snake_case 键 )。
type flatPayload struct {
	BusinessID string `json:"business_id"`
	UserID     uint64 `json:"user_id"`
	SessionID  string `json:"session_id"`
}

// testAudience 为 FlatClaims 单元测试使用的受众。
const testAudience = "aud-1"

func TestFlatClaims_Marshal(t *testing.T) {
	t.Parallel()

	fc := &FlatClaims[flatPayload]{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    defaultIssuer,
			Audience:  jwt.ClaimStrings{testAudience},
			ExpiresAt: jwt.NewNumericDate(time.Unix(1700000000, 0)),
		},
		Data: flatPayload{BusinessID: "biz-1", UserID: 1, SessionID: "sess-1"},
	}

	b, err := json.Marshal(fc)
	require.NoError(t, err)

	// 序列化结果为单一顶层 object:
	// 注册 claims 与业务 claims 平铺，无 "Data" 包装键。
	var m map[string]any
	require.NoError(t, json.Unmarshal(b, &m))

	assert.Equal(t, defaultIssuer, m["iss"])
	assert.Equal(t, []any{testAudience}, m["aud"])
	assert.Equal(t, float64(1700000000), m["exp"])
	assert.Equal(t, "biz-1", m["business_id"])
	assert.Equal(t, float64(1), m["user_id"])
	assert.Equal(t, "sess-1", m["session_id"])

	_, hasDataKey := m["Data"]
	assert.False(t, hasDataKey)
}

func TestFlatClaims_Unmarshal(t *testing.T) {
	t.Parallel()

	raw := []byte(`{"iss":"jit","aud":"aud-1","exp":1700000000,"business_id":"biz-1","user_id":1,"session_id":"sess-1"}`)

	fc := &FlatClaims[flatPayload]{}
	require.NoError(t, json.Unmarshal(raw, fc))

	// 保留键归入 RegisteredClaims，其余归入 Data。
	assert.Equal(t, defaultIssuer, fc.Issuer)
	assert.Equal(t, jwt.ClaimStrings{testAudience}, fc.Audience)
	require.NotNil(t, fc.ExpiresAt)
	assert.Equal(t, int64(1700000000), fc.ExpiresAt.Unix())
	assert.Equal(t, flatPayload{BusinessID: "biz-1", UserID: 1, SessionID: "sess-1"}, fc.Data)
}

// TestFlatClaims_ReservedKeyConflict 业务数据携带保留键时必须拒绝签发，
// 避免静默覆盖安全相关的注册 claims。
func TestFlatClaims_ReservedKeyConflict(t *testing.T) {
	t.Parallel()

	type evilPayload struct {
		Iss string `json:"iss"`
	}

	fc := &FlatClaims[evilPayload]{
		RegisteredClaims: jwt.RegisteredClaims{Issuer: defaultIssuer},
		Data:             evilPayload{Iss: "attacker"},
	}

	_, err := json.Marshal(fc)
	assert.Error(t, err)
}

// TestFlatClaims_NonObjectData 非对象类型无法平铺，必须报错。
func TestFlatClaims_NonObjectData(t *testing.T) {
	t.Parallel()

	fc := &FlatClaims[string]{Data: "not-an-object"}

	_, err := json.Marshal(fc)
	assert.Error(t, err)
}
