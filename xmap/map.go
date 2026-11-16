package xmap

import "github.com/jrmarcco/jit/internal/errs"

// Keys 返回 map 中所有 key 组成的切片。
// 基于迭代器的替代写法可参考 slices.Collect(maps.Keys(m))。
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Vals 返回 map 中所有 value 组成的切片。
// 基于迭代器的替代写法可参考 slices.Collect(maps.Values(m))。
func Vals[K comparable, V any](m map[K]V) []V {
	vals := make([]V, 0, len(m))
	for _, v := range m {
		vals = append(vals, v)
	}
	return vals
}

// MapKV 为 map 的键值对。
type MapKV[K comparable, V any] struct {
	Key K
	Val V
}

// KeysVals 返回 map 中所有键值对组成的切片。
func KeysVals[K comparable, V any](m map[K]V) []MapKV[K, V] {
	keys := make([]MapKV[K, V], 0, len(m))
	for k, v := range m {
		keys = append(keys, MapKV[K, V]{Key: k, Val: v})
	}
	return keys
}

// Merge 合并多个 map，重复 key 由后出现的 map 覆盖。
func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
	defaultMergeFunc := func(_ V, second V) V {
		return second
	}
	return MergeFunc(defaultMergeFunc, maps...)
}

// MergeFunc 合并多个 map，重复 key 的取值由 mergeFunc 决定
// ( 入参为已存在的值与新值，返回最终保留的值 )。
func MergeFunc[K comparable, V any](mergeFunc func(_ V, second V) V, maps ...map[K]V) map[K]V {
	totalLen := 0
	for _, m := range maps {
		totalLen += len(m)
	}

	res := make(map[K]V, totalLen)
	for _, m := range maps {
		for k, v := range m {
			if val, ok := res[k]; ok {
				res[k] = mergeFunc(val, v)
				continue
			}
			res[k] = v
		}
	}
	return res
}

// ToMap 由等长的 keys 与 vals 切片构建 map。
// nil 切片视为空切片；长度不一致时返回错误；
// 重复 key 保留最后一个值。
func ToMap[K comparable, V any](keys []K, vals []V) (map[K]V, error) {
	keyLen := len(keys)
	if keyLen != len(vals) {
		return nil, errs.NewErrInvalidKeyValLen()
	}

	res := make(map[K]V, keyLen)
	for i, k := range keys {
		res[k] = vals[i]
	}
	return res, nil
}
