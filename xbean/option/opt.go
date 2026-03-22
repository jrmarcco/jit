package option

// Opt is generalized design for option pattern.
type Opt[T any] func(*T)

// Apply applies the options to the target.
func Apply[T any](t *T, opts ...Opt[T]) {
	for _, opt := range opts {
		opt(t)
	}
}

// OptErr is a variant of Opt that returns an error.
type OptErr[T any] func(*T) error

// ApplyErr applies the options to the target and returns an error if any option returns an error.
func ApplyErr[T any](t *T, opts ...OptErr[T]) error {
	for _, opt := range opts {
		if err := opt(t); err != nil {
			return err
		}
	}
	return nil
}
