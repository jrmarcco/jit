package xjwt

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// FlatClaims 把泛型业务数据 T 平铺到 jwt payload 顶层，
// 与 CustomClaims[T] ( 业务数据序列化在 "Data" 键下 ) 相对。
//
// 适用场景: 网关 / 第三方按顶层键直接解析业务 claims，例如:
//
//	{ "iss": "...", "aud": "...", "business_id": "...", "user_id": 1 }
//
// T 的 JSON 键由调用方的 struct tag 决定。
// 约束:
//   - T 必须序列化为 JSON object ( struct / map )；
//   - T 的键不得与 RFC 7519 保留键 ( iss/sub/aud/exp/nbf/iat/jti ) 冲突，
//     冲突时 Marshal 返回错误。
type FlatClaims[T any] struct {
	jwt.RegisteredClaims
	Data T
}

var _ jwt.Claims = (*FlatClaims[any])(nil)

// MarshalJSON 合并 RegisteredClaims 与 Data 到同一顶层 JSON object。
func (fc *FlatClaims[T]) MarshalJSON() ([]byte, error) {
	dataRaw, err := json.Marshal(fc.Data)
	if err != nil {
		return nil, fmt.Errorf("[jit] failed to marshal flat claims data: %w", err)
	}

	dataFields := make(map[string]json.RawMessage)
	if err = json.Unmarshal(dataRaw, &dataFields); err != nil {
		return nil, errors.New("[jit] flat claims data must be a json object")
	}

	regRaw, err := json.Marshal(fc.RegisteredClaims)
	if err != nil {
		return nil, fmt.Errorf("[jit] failed to marshal registered claims: %w", err)
	}

	regFields := make(map[string]json.RawMessage)
	if err = json.Unmarshal(regRaw, &regFields); err != nil {
		return nil, fmt.Errorf("[jit] failed to marshal registered claims: %w", err)
	}

	for key := range regFields {
		if _, ok := dataFields[key]; ok {
			return nil, fmt.Errorf("[jit] flat claims data must not contain reserved claim key %q", key)
		}
	}

	merged := make(map[string]json.RawMessage, len(regFields)+len(dataFields))
	for key, val := range regFields {
		merged[key] = val
	}
	for key, val := range dataFields {
		merged[key] = val
	}

	return json.Marshal(merged)
}

// UnmarshalJSON 把顶层键拆分回 RegisteredClaims 与 Data:
// RFC 7519 保留键归入 RegisteredClaims，其余归入 Data。
func (fc *FlatClaims[T]) UnmarshalJSON(b []byte) error {
	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(b, &fields); err != nil {
		return fmt.Errorf("[jit] failed to unmarshal flat claims: %w", err)
	}

	regFields := make(map[string]json.RawMessage)
	dataFields := make(map[string]json.RawMessage)
	for key, val := range fields {
		if isRegisteredClaimKey(key) {
			regFields[key] = val
		} else {
			dataFields[key] = val
		}
	}

	if len(regFields) > 0 {
		regRaw, err := json.Marshal(regFields)
		if err != nil {
			return fmt.Errorf("[jit] failed to unmarshal registered claims: %w", err)
		}
		if err := json.Unmarshal(regRaw, &fc.RegisteredClaims); err != nil {
			return fmt.Errorf("[jit] failed to unmarshal registered claims: %w", err)
		}
	}

	if len(dataFields) > 0 {
		dataRaw, err := json.Marshal(dataFields)
		if err != nil {
			return fmt.Errorf("[jit] failed to unmarshal flat claims data: %w", err)
		}
		if err := json.Unmarshal(dataRaw, &fc.Data); err != nil {
			return fmt.Errorf("[jit] failed to unmarshal flat claims data: %w", err)
		}
	}

	return nil
}

func isRegisteredClaimKey(key string) bool {
	switch key {
	case "iss", "sub", "aud", "exp", "nbf", "iat", "jti":
		return true
	default:
		return false
	}
}
