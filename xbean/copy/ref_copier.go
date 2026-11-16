package copier

import (
	"maps"
	"reflect"
	"slices"
	"time"

	"github.com/jrmarcco/jit/xbean/option"
	"github.com/jrmarcco/jit/xset"
)

// RefCopier 为基于反射的 struct 拷贝器。
// 创建时预解析源类型与目标类型的字段映射树，拷贝时按树遍历，避免重复反射。
type RefCopier[S any, D any] struct {
	root        fieldNode
	atomicTypes []reflect.Type

	defaultConf copyConf
}

// NewRefCopier 创建 RefCopier。
// S 与 D 必须均为 struct 类型，否则返回错误；
// opts 在创建期生效并作为后续拷贝的默认配置。
func NewRefCopier[S any, D any](opts ...option.Opt[copyConf]) (*RefCopier[S, D], error) {
	srcTyp := reflect.TypeFor[S]()
	dstTyp := reflect.TypeFor[D]()

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

// createFieldNode 递归构建 srcTyp 与 dstTyp 的字段映射树:
// 同名且类型兼容的导出字段生成一个节点。
// 字段为指针的指针时不支持，返回错误；
// 非内建、非原子、非 struct 的字段 ( 如 interface、func ) 跳过。
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
				// 指针的指针。
				return errPtrToPtr(srcFd.Name)
			}

			if dstFd.Type.Kind() == reflect.Pointer && dstFd.Type.Elem().Kind() == reflect.Pointer {
				// 指针的指针。
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
				// 内建类型，节点为叶子节点。
			case rc.isAtomicType(srcFdTyp):
				// 原子类型 ( 如 time.Time )，节点为叶子节点。
			case srcFdTyp.Kind() == reflect.Struct:
				// struct 类型，递归下钻。
				if err := rc.createFieldNode(srcFdTyp, dstFdTyp, &node); err != nil {
					return err
				}
			default:
				// 非内建、非原子、非 struct，无法拷贝，跳过。
				continue
			}

			// 把节点加入 root.fields。
			root.fields = append(root.fields, node)
		}
	}
	return nil
}

// isAtomicType 判断类型是否为原子类型 ( 整体拷贝、不下钻 )。
func (rc *RefCopier[S, D]) isAtomicType(typ reflect.Type) bool {
	return slices.Contains(rc.atomicTypes, typ)
}

// Copy 新建 D 并从 src 拷贝匹配字段。
//
// 拷贝语义见 CopyTo:
//   - 源字段为零值时跳过 ( 目标保留原值 )；
//   - slice/map 字段为浅拷贝 ( 共享底层数据 )。
func (rc *RefCopier[S, D]) Copy(src *S, opts ...option.Opt[copyConf]) (*D, error) {
	dst := new(D)
	err := rc.CopyTo(src, dst, opts...)
	return dst, err
}

// CopyTo 把 src 中匹配的导出字段拷贝到 dst。
//
// 拷贝语义:
//   - 仅拷贝字段名相同且类型兼容的导出字段。
//   - 源字段为零值时跳过: 目标字段保留原值
//     ( 对已填充的 dst 调用 CopyTo 不会清空字段 )。
//   - 源指针字段为 nil 时同样跳过。
//   - slice/map/chan 字段为浅拷贝: 源与目标共享底层数据。
//   - 经 ConvertFd 注册转换器的字段走转换器。
//
// 无任何匹配字段时 CopyTo 返回 nil ( 无事可做 )。
func (rc *RefCopier[S, D]) CopyTo(src *S, dst *D, opts ...option.Opt[copyConf]) error {
	if len(rc.root.fields) == 0 {
		return nil
	}

	cc := rc.defaultCopyConf()
	option.Apply(&cc, opts...)
	return rc.copyTree(src, dst, cc)
}

// defaultCopyConf 深拷贝默认配置，
// 避免单次拷贝的临时配置污染创建期配置。
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
		cc.covertFds = make(map[string]convertor, len(rc.defaultConf.covertFds))
	}

	maps.Copy(cc.covertFds, rc.defaultConf.covertFds)

	return cc
}

// copyTree 从根节点开始拷贝。
func (rc *RefCopier[S, D]) copyTree(src *S, dst *D, cc copyConf) error {
	srcTyp := reflect.TypeOf(src)
	srcVal := reflect.ValueOf(src)

	dstTyp := reflect.TypeOf(dst)
	dstVal := reflect.ValueOf(dst)

	return rc.copyNode(srcTyp, srcVal, dstTyp, dstVal, &rc.root, cc)
}

// copyNode 拷贝单个节点:
// 叶子节点 ( 无子字段 ) 走 copyLeafNode，struct 节点走 copyStructNode。
// 源指针为 nil 时直接跳过。
func (rc *RefCopier[S, D]) copyNode(srcTyp reflect.Type, srcVal reflect.Value, dstTyp reflect.Type, dstVal reflect.Value, root *fieldNode, cc copyConf) error {
	// 保留解引用前的原值，供字段转换器使用。
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

// derefSrc 解引用源值，nil 指针返回 false ( 跳过拷贝 )。
func (rc *RefCopier[S, D]) derefSrc(typ reflect.Type, val reflect.Value) (reflect.Type, reflect.Value, bool) {
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return nil, reflect.Value{}, false
		}
		return typ.Elem(), val.Elem(), true
	}
	return typ, val, true
}

// derefDst 解引用目标值，nil 指针先分配再解引用。
func (rc *RefCopier[S, D]) derefDst(typ reflect.Type, val reflect.Value) (reflect.Type, reflect.Value) {
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			val.Set(reflect.New(typ.Elem()))
		}
		return typ.Elem(), val.Elem()
	}
	return typ, val
}

// copyLeafNode 拷贝叶子字段:
//   - 注册了转换器: 经转换器转换后写入原目标值 ( 可能为指针字段 )。
//   - 未注册转换器: 类型一致且源值非零时直接赋值。
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

// copyStructNode 遍历映射树的子字段并逐个拷贝。
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

// fieldNode 为字段映射树的节点。
type fieldNode struct {
	name   string
	fields []fieldNode
	sIndex int // 源字段下标
	dIndex int // 目标字段下标
}

// isBuiltinType 判断类型种类是否可直接赋值拷贝
// ( 含 slice/map/chan/array，均为浅拷贝共享底层 )。
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
