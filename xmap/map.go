package xmap

import "github.com/jrmarcco/jit/internal/errs"

// Keys return a slice of the keys in the map.
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Vals returns a slice of the values in the map.
func Vals[K comparable, V any](m map[K]V) []V {
	vals := make([]V, 0, len(m))
	for _, v := range m {
		vals = append(vals, v)
	}
	return vals
}

// MapKV is a key-value pair of a map.
type MapKV[K comparable, V any] struct {
	Key K
	Val V
}

// KeysVals returns a slice of the key-value pairs in the map.
func KeysVals[K comparable, V any](m map[K]V) []MapKV[K, V] {
	keys := make([]MapKV[K, V], 0, len(m))
	for k, v := range m {
		keys = append(keys, MapKV[K, V]{Key: k, Val: v})
	}
	return keys
}

func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
	defaultMergeFunc := func(_ V, second V) V {
		return second
	}
	return MergeFunc(defaultMergeFunc, maps...)
}

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

// ToMap converts a slice of keys and a slice of values to a map.
func ToMap[K comparable, V any](keys []K, vals []V) (map[K]V, error) {
	if keys == nil || vals == nil {
		return nil, errs.NilErr("keys or vals")
	}

	keyLen := len(keys)
	if keyLen != len(vals) {
		return nil, errs.ErrInvalidKeyValLen()
	}

	res := make(map[K]V, keyLen)
	for i, k := range keys {
		res[k] = vals[i]
	}
	return res, nil
}
