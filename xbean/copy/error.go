package copier

import (
	"errors"
	"fmt"
	"reflect"
)

var errConvertFdTypeMismatch = errors.New("[jit] convert func type mismatch")

func errInvalidType(want, got any) error {
	return fmt.Errorf("[jit] invalid type: want %s, got %#v", want, got)
}

func errPtrToPtr(name string) error {
	return fmt.Errorf("[jit] cannot copy pointer to pointer: %s", name)
}

func errFieldTypeMismatch(name string, srcTyp, dstTyp reflect.Type) error {
	return fmt.Errorf("[jit] type mismatch at field %s: source type %s != destination type %s", name, srcTyp, dstTyp)
}
