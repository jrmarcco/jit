package copier

import (
	"errors"
	"fmt"
	"reflect"
)

// errConvertFdTypeMismatch 为字段转换器收到的源值类型不匹配时返回的错误。
var errConvertFdTypeMismatch = errors.New("[jit] convert func type mismatch")

// errInvalidType 返回类型不符合预期的错误。
func errInvalidType(want, got any) error {
	return fmt.Errorf("[jit] invalid type: want %s, got %#v", want, got)
}

// errPtrToPtr 返回字段为指针的指针时不支持的错误。
func errPtrToPtr(name string) error {
	return fmt.Errorf("[jit] cannot copy pointer to pointer: %s", name)
}

// errFieldTypeMismatch 返回字段类型不匹配的错误。
func errFieldTypeMismatch(name string, srcTyp, dstTyp reflect.Type) error {
	return fmt.Errorf("[jit] type mismatch at field %s: source type %s != destination type %s", name, srcTyp, dstTyp)
}
