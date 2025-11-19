package copier

import (
	"testing"

	"github.com/JrMarcco/jit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefCopier_Copy(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		copyFn  func() (any, error)
		wantDst any
		wantErr error
	}{
		{
			name: "struct basic",
			copyFn: func() (any, error) {
				copier, err := NewRefCopier[basicSrc, basicDst]()
				if err != nil {
					return nil, err
				}

				return copier.Copy(&basicSrc{
					IntVal: 1,
					IntPtr: jit.Ptr(1),
					StrVal: "test",
				})
			},
			wantDst: &basicDst{
				IntVal: 1,
				IntPtr: jit.Ptr(1),
				StrVal: "test",
			},
		}, {
			name: "struct with no field",
			copyFn: func() (any, error) {
				copier, err := NewRefCopier[noFdSrc, noFdDst]()
				if err != nil {
					return nil, err
				}

				return copier.Copy(&noFdSrc{})
			},
			wantDst: &noFdDst{},
		}, {
			name: "struct with private field",
			copyFn: func() (any, error) {
				copier, err := NewRefCopier[priFdSrc, priFdDst]()
				if err != nil {
					return nil, err
				}

				return copier.Copy(&priFdSrc{
					IntVal: 1,
					StrVal: "test",
					pri:    param{Val: "pri"},
					priPtr: &param{Val: "priPtr"},
				})
			},
			wantDst: &priFdDst{
				IntVal: 1,
				StrVal: "test",
				priVal: param{},
				priPtr: nil,
			},
		}, {
			name: "struct with slice",
			copyFn: func() (any, error) {
				copier, err := NewRefCopier[sliceSrc, sliceDst]()
				if err != nil {
					return nil, err
				}

				return copier.Copy(&sliceSrc{
					Ints:   []int{1, 2, 3},
					Strs:   []string{"a", "b", "c"},
					Params: []param{{Val: "a"}, {Val: "b"}, {Val: "c"}},
				})
			},
			wantDst: &sliceDst{
				Ints:   []int{1, 2, 3},
				Strs:   []string{"a", "b", "c"},
				Params: []param{{Val: "a"}, {Val: "b"}, {Val: "c"}},
			},
		}, {
			name: "struct complex",
			copyFn: func() (any, error) {
				copier, err := NewRefCopier[complexSrc, complexDst]()
				if err != nil {
					return nil, err
				}

				return copier.Copy(&complexSrc{
					Basic: basicSrc{
						IntVal: 1,
						IntPtr: jit.Ptr(1),
						StrVal: "test",
					},
					Embed: embedSrc{
						Basic: basicSrc{
							IntVal: 1,
							IntPtr: jit.Ptr(1),
							StrVal: "test",
						},
						Slice: &sliceSrc{
							Ints:   []int{1, 2, 3},
							Strs:   []string{"a", "b", "c"},
							Params: []param{{Val: "a"}, {Val: "b"}, {Val: "c"}},
						},
					},
				})
			},
			wantDst: &complexDst{
				Basic: basicDst{
					IntVal: 1,
					IntPtr: jit.Ptr(1),
					StrVal: "test",
				},
				Embed: embedDst{
					Basic: basicDst{
						IntVal: 1,
						IntPtr: jit.Ptr(1),
						StrVal: "test",
					},
					Slice: &sliceDst{
						Ints:   []int{1, 2, 3},
						Strs:   []string{"a", "b", "c"},
						Params: []param{{Val: "a"}, {Val: "b"}, {Val: "c"}},
					},
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dst, err := tc.copyFn()
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}

			assert.Equal(t, tc.wantDst, dst)
		})
	}
}

type param struct {
	Val string
}

type basicSrc struct {
	IntVal int
	IntPtr *int
	StrVal string
}

type basicDst struct {
	IntVal int
	IntPtr *int
	StrVal string
}

type noFdSrc struct{}

type noFdDst struct{}

type priFdSrc struct {
	IntVal int
	StrVal string
	pri    param
	priPtr *param
}

type priFdDst struct {
	IntVal int
	StrVal string
	priVal param
	priPtr *param
}

type sliceSrc struct {
	Ints   []int
	Strs   []string
	Params []param
}

type sliceDst struct {
	Ints   []int
	Strs   []string
	Params []param
}

type embedSrc struct {
	Basic basicSrc
	Slice *sliceSrc
}

type embedDst struct {
	Basic basicDst
	Slice *sliceDst
}

type complexSrc struct {
	Basic basicSrc
	Embed embedSrc
}

type complexDst struct {
	Basic basicDst
	Embed embedDst
}

func BenchmarkRefCopier_Copy(b *testing.B) {
	b.Run("reuse", func(b *testing.B) {
		copier, err := NewRefCopier[complexSrc, complexDst]()
		require.NoError(b, err)

		for b.Loop() {
			_, err := copier.Copy(&complexSrc{
				Basic: basicSrc{
					IntVal: 1,
					IntPtr: jit.Ptr(1),
					StrVal: "test",
				},
				Embed: embedSrc{
					Basic: basicSrc{
						IntVal: 1,
						IntPtr: jit.Ptr(1),
						StrVal: "test",
					},
					Slice: &sliceSrc{
						Ints:   []int{1, 2, 3},
						Strs:   []string{"a", "b", "c"},
						Params: []param{{Val: "a"}, {Val: "b"}, {Val: "c"}},
					},
				},
			})
			require.NoError(b, err)
		}
	})
}
