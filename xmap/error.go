package xmap

import "errors"

// ErrNilComparator 为比较器为 nil 时返回的错误。
var ErrNilComparator = errors.New("[jit] comparator can not be nil")
