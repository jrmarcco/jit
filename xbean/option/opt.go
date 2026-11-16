// Package option 提供泛型化的 option 模式:
// Opt[T] 修改目标，OptErr[T] 在此基础上返回错误。
package option

// Opt 为 option 模式的泛型定义。
type Opt[T any] func(*T)

// Apply 依次把 opts 应用到 t。
func Apply[T any](t *T, opts ...Opt[T]) {
	for _, opt := range opts {
		opt(t)
	}
}

// OptErr 为返回 error 的 Opt 变体。
type OptErr[T any] func(*T) error

// ApplyErr 依次把 opts 应用到 t，任一 opt 返回错误时立即返回该错误。
func ApplyErr[T any](t *T, opts ...OptErr[T]) error {
	for _, opt := range opts {
		if err := opt(t); err != nil {
			return err
		}
	}
	return nil
}
