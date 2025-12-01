package copier

import (
	"github.com/jrmarcco/jit/bean/copy/converter"
	"github.com/jrmarcco/jit/bean/option"
	"github.com/jrmarcco/jit/xset"
)

// Copier is a type that can copy a source object to a destination object.
type Copier[S any, D any] interface {
	Copy(src *S) (*D, error)
	CopyTo(src *S, dst *D) error
}

type convertFunc func(src any) (any, error)

type copyConf struct {
	ignoreFds *xset.MapSet[string]
	covertFds map[string]convertFunc
}

func newCopyConf() copyConf {
	return copyConf{}
}

func (cc *copyConf) InIgnore(fd string) bool {
	if cc.ignoreFds == nil {
		return false
	}
	return cc.ignoreFds.Exist(fd)
}

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

func ConvertFd[S any, D any](fd string, converter converter.Converter[S, D]) option.Opt[copyConf] {
	return func(cc *copyConf) {
		if fd == "" || converter == nil {
			return
		}

		if cc.covertFds == nil {
			const defaultSize = 8
			cc.covertFds = make(map[string]convertFunc, defaultSize)
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
