// Package copier 提供基于反射的 struct 拷贝器，
// 支持字段忽略与字段转换。
package copier

import (
	"github.com/jrmarcco/jit/xbean/copy/convert"
	"github.com/jrmarcco/jit/xbean/option"
	"github.com/jrmarcco/jit/xset"
)

// Copier 为把源对象拷贝到目标对象的抽象。
type Copier[S any, D any] interface {
	// Copy 新建 D 并执行拷贝。
	Copy(src *S) (*D, error)

	// CopyTo 把 src 拷贝到已存在的 dst。
	CopyTo(src *S, dst *D) error
}

// convertor 为字段转换函数的统一签名。
type convertor func(src any) (any, error)

// copyConf 为拷贝配置 ( 忽略字段与字段转换器 )。
type copyConf struct {
	ignoreFds *xset.MapSet[string]
	covertFds map[string]convertor
}

// newCopyConf 创建空配置。
func newCopyConf() copyConf {
	return copyConf{}
}

// InIgnore 判断字段是否在忽略列表中。
func (cc *copyConf) InIgnore(fd string) bool {
	if cc.ignoreFds == nil {
		return false
	}
	return cc.ignoreFds.Exist(fd)
}

// IgnoreFds 配置需要忽略的字段列表。
func IgnoreFds(fds ...string) option.Opt[copyConf] {
	return func(cc *copyConf) {
		if len(fds) == 0 {
			return
		}

		if cc.ignoreFds == nil {
			cc.ignoreFds = xset.NewMapSet[string](len(fds))
		}

		for _, fd := range fds {
			cc.ignoreFds.Add(fd)
		}
	}
}

// ConvertFd 为字段 fd 注册类型转换器 converter。
// converter 为 nil 或 fd 为空串时忽略该配置。
func ConvertFd[S any, D any](fd string, converter convert.Convertor[S, D]) option.Opt[copyConf] {
	return func(cc *copyConf) {
		if fd == "" || converter == nil {
			return
		}

		if cc.covertFds == nil {
			const defaultSize = 8
			cc.covertFds = make(map[string]convertor, defaultSize)
		}

		cc.covertFds[fd] = func(src any) (any, error) {
			var dst D
			srcVal, ok := src.(S)
			if !ok {
				return dst, errConvertFdTypeMismatch
			}

			return converter.Convert(srcVal)
		}
	}
}
