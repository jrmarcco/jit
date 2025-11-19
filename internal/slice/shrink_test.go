package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShrink(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		originCap int
		originLen int
		wantCap   int
	}{
		{
			name:      "empty slice",
			originCap: 0,
			originLen: 0,
			wantCap:   0,
		}, {
			name:      "huge capacity: when the ratio >= 2, shrink to 1.5 times of the original capacity",
			originCap: 8192,
			originLen: 1024,
			wantCap:   1536,
		}, {
			name:      "large capacity: when the ratio >= 2, shrink to 50% of the original capacity",
			originCap: 2048,
			originLen: 256,
			wantCap:   1024,
		}, {
			name:      "medium capacity: when the ratio >= 2.5, shrink to 62.5% of the original capacity",
			originCap: 1024,
			originLen: 256,
			wantCap:   640,
		}, {
			name:      "small capacity: when the ratio >= 3, shrink to 50% of the original capacity",
			originCap: 128,
			originLen: 8,
			wantCap:   64,
		}, {
			name:      "small capacity: when the ratio < 3, not shrink",
			originCap: 128,
			originLen: 64,
			wantCap:   128,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			slice := make([]int, 0, tc.originCap)
			for i := 0; i < tc.originLen; i++ {
				slice = append(slice, i)
			}

			res := Shrink(slice)
			assert.Equal(t, tc.wantCap, cap(res))
		})
	}
}
