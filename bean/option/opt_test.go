package option

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type example struct {
	Int int
	Str string
}

func TestApply(t *testing.T) {
	t.Parallel()

	e := &example{}
	Apply(e, func(e *example) {
		e.Int = 1
		e.Str = "test"
	})

	assert.Equal(t, e.Int, 1)
	assert.Equal(t, e.Str, "test")
}

func withStrErr(str string) OptErr[example] {
	return func(e *example) error {
		if str == "" {
			return errors.New("str error")
		}
		e.Str = str
		return nil
	}
}

func TestApplyErr(t *testing.T) {
	t.Parallel()

	e := &example{}
	err := ApplyErr(e, func(e *example) error {
		e.Int = 1
		e.Str = "test"
		return nil
	})
	require.NoError(t, err)
	assert.Equal(t, e.Int, 1)
	assert.Equal(t, e.Str, "test")

	err = ApplyErr(e, withStrErr("test"))
	require.NoError(t, err)
	assert.Equal(t, e.Str, "test")

	err = ApplyErr(e, withStrErr(""))
	assert.Error(t, err, errors.New("str error"))
	assert.Equal(t, e.Int, 1)
	assert.Equal(t, e.Str, "test")
}
