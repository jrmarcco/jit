package copier

import (
	"maps"
	"reflect"
	"slices"
	"time"

	"github.com/JrMarcco/jit/bean/option"
	"github.com/JrMarcco/jit/xset"
)

type RefCopier[S any, D any] struct {
	root        fieldNode
	atomicTypes []reflect.Type

	defaultConf copyConf
}

func (rc *RefCopier[S, D]) createFieldNode(srcTyp, dstTyp reflect.Type, root *fieldNode) error {
	srcMap := map[string]int{}
	for i := range srcTyp.NumField() {
		fd := srcTyp.Field(i)
		if fd.IsExported() {
			srcMap[fd.Name] = i
		}
	}

	for dstIdx := range dstTyp.NumField() {
		dstFd := dstTyp.Field(dstIdx)

		if !dstFd.IsExported() {
			continue
		}

		if srcIdx, ok := srcMap[dstFd.Name]; ok {
			srcFd := srcTyp.Field(srcIdx)

			if srcFd.Type.Kind() == reflect.Pointer && srcFd.Type.Elem().Kind() == reflect.Pointer {
				// pointer to pointer.
				return errPtrToPtr(srcFd.Name)
			}

			if dstFd.Type.Kind() == reflect.Pointer && dstFd.Type.Elem().Kind() == reflect.Pointer {
				// pointer to pointer.
				return errPtrToPtr(dstFd.Name)
			}

			node := fieldNode{
				name:   dstFd.Name,
				sIndex: srcIdx,
				dIndex: dstIdx,
				fields: []fieldNode{},
			}

			srcFdTyp := srcFd.Type
			dstFdTyp := dstFd.Type

			if srcFdTyp.Kind() == reflect.Pointer {
				srcFdTyp = srcFdTyp.Elem()
			}

			if dstFdTyp.Kind() == reflect.Pointer {
				dstFdTyp = dstFdTyp.Elem()
			}

			switch {
			case isBuiltinType(srcFdTyp.Kind()):
				// builtin type, node is leaf node.
			case rc.isAtomicType(srcFdTyp):
				// is atomic type, node is leaf node.
			case srcFdTyp.Kind() == reflect.Struct:
				// is struct type.
				if err := rc.createFieldNode(srcFdTyp, dstFdTyp, &node); err != nil {
					return err
				}
			default:
				// is not builtin type, not struct type, not atomic type.
				// can not copy it, skip.
				continue
			}

			// add node to root.fields.
			root.fields = append(root.fields, node)
		}
	}
	return nil
}

func (rc *RefCopier[S, D]) isAtomicType(typ reflect.Type) bool {
	return slices.Contains(rc.atomicTypes, typ)
}

func (rc *RefCopier[S, D]) Copy(src *S, opts ...option.Opt[copyConf]) (*D, error) {
	dst := new(D)
	err := rc.CopyTo(src, dst, opts...)
	return dst, err
}

func (rc *RefCopier[S, D]) CopyTo(src *S, dst *D, opts ...option.Opt[copyConf]) error {
	if len(rc.root.fields) == 0 {
		return nil
	}

	cc := rc.defaultCopyConf()
	option.Apply(&cc, opts...)
	return rc.copyTree(src, dst, cc)
}

func (rc *RefCopier[S, D]) defaultCopyConf() copyConf {
	cc := newCopyConf()

	if rc.defaultConf.ignoreFds != nil {
		ignoreFds := xset.NewMapSet[string](rc.defaultConf.ignoreFds.Size())

		for _, fd := range rc.defaultConf.ignoreFds.Elems() {
			ignoreFds.Add(fd)
		}

		cc.ignoreFds = ignoreFds
	}

	if cc.covertFds == nil {
		cc.covertFds = make(map[string]convertFunc, len(rc.defaultConf.covertFds))
	}

	maps.Copy(cc.covertFds, rc.defaultConf.covertFds)

	return cc
}

func (rc *RefCopier[S, D]) copyTree(src *S, dst *D, cc copyConf) error {
	srcTyp := reflect.TypeOf(src)
	srcVal := reflect.ValueOf(src)

	dstTyp := reflect.TypeOf(dst)
	dstVal := reflect.ValueOf(dst)

	return rc.copyNode(srcTyp, srcVal, dstTyp, dstVal, &rc.root, cc)
}

func (rc *RefCopier[S, D]) copyNode(srcTyp reflect.Type, srcVal reflect.Value, dstTyp reflect.Type, dstVal reflect.Value, root *fieldNode, cc copyConf) error {
	oriSrcVal := srcVal
	oriDstVal := dstVal

	var ok bool
	srcTyp, srcVal, ok = rc.derefSrc(srcTyp, srcVal)
	if !ok {
		return nil
	}

	dstTyp, dstVal = rc.derefDst(dstTyp, dstVal)
	if len(root.fields) == 0 {
		return rc.copyLeafNode(srcTyp, srcVal, dstTyp, dstVal, oriSrcVal, oriDstVal, root, cc)
	}

	return rc.copyStructNode(srcTyp, srcVal, dstTyp, dstVal, root, cc)
}

func (rc *RefCopier[S, D]) derefSrc(typ reflect.Type, val reflect.Value) (reflect.Type, reflect.Value, bool) {
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return nil, reflect.Value{}, false
		}
		return typ.Elem(), val.Elem(), true
	}
	return typ, val, true
}

func (rc *RefCopier[S, D]) derefDst(typ reflect.Type, val reflect.Value) (reflect.Type, reflect.Value) {
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			val.Set(reflect.New(typ.Elem()))
		}
		return typ.Elem(), val.Elem()
	}
	return typ, val
}

func (rc *RefCopier[S, D]) copyLeafNode(
	srcTyp reflect.Type, srcVal reflect.Value,
	dstTyp reflect.Type, dstVal reflect.Value,
	oriSrcVal, oriDstVal reflect.Value,
	root *fieldNode,
	cc copyConf,
) error {
	fdName := root.name
	if !dstVal.CanSet() {
		return nil
	}

	convertFunc, ok := cc.covertFds[fdName]
	if !ok {
		if srcTyp != dstTyp {
			return errFieldTypeMismatch(fdName, srcTyp, dstTyp)
		}
		if srcVal.IsZero() {
			return nil
		}
		dstVal.Set(srcVal)
		return nil
	}

	if !oriDstVal.CanSet() {
		return nil
	}

	srcConverted, err := convertFunc(oriSrcVal.Interface())
	if err != nil {
		return err
	}

	srcConvTyp := reflect.TypeOf(srcConverted)
	srcConvVal := reflect.ValueOf(srcConverted)
	if srcConvTyp != oriDstVal.Type() {
		return errFieldTypeMismatch(fdName, srcConvTyp, oriDstVal.Type())
	}

	oriDstVal.Set(srcConvVal)
	return nil
}

func (rc *RefCopier[S, D]) copyStructNode(srcTyp reflect.Type, srcVal reflect.Value, dstTyp reflect.Type, dstVal reflect.Value, root *fieldNode, cc copyConf) error {
	for _, field := range root.fields {
		if cc.InIgnore(field.name) {
			continue
		}

		srcFdTyp := srcTyp.Field(field.sIndex)
		srcFdVal := srcVal.Field(field.sIndex)

		dstFdTyp := dstTyp.Field(field.dIndex)
		dstFdVal := dstVal.Field(field.dIndex)

		if err := rc.copyNode(srcFdTyp.Type, srcFdVal, dstFdTyp.Type, dstFdVal, &field, cc); err != nil {
			return err
		}
	}
	return nil
}

func NewRefCopier[S any, D any](opts ...option.Opt[copyConf]) (*RefCopier[S, D], error) {
	srcTyp := reflect.TypeOf(new(S)).Elem()
	dstTyp := reflect.TypeOf(new(D)).Elem()

	if srcTyp.Kind() != reflect.Struct {
		return nil, errInvalidType("struct", srcTyp)
	}

	if dstTyp.Kind() != reflect.Struct {
		return nil, errInvalidType("struct", dstTyp)
	}

	root := fieldNode{
		fields: []fieldNode{},
	}

	copier := &RefCopier[S, D]{
		root: root,
		atomicTypes: []reflect.Type{
			reflect.TypeOf(time.Time{}),
		},
	}

	if err := copier.createFieldNode(srcTyp, dstTyp, &root); err != nil {
		return nil, err
	}

	copier.root = root

	cc := newCopyConf()
	option.Apply(&cc, opts...)

	copier.defaultConf = cc
	return copier, nil
}

type fieldNode struct {
	name   string
	fields []fieldNode
	sIndex int // source index
	dIndex int // destination index
}

func isBuiltinType(kind reflect.Kind) bool {
	switch kind {
	case
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128,
		reflect.String,
		reflect.Slice,
		reflect.Map,
		reflect.Chan,
		reflect.Array:
		return true
	default:
		return false
	}
}
