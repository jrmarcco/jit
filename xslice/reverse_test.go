package xslice

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverseOfInt(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		slice []int
		want  []int
	}{
		{
			name:  tcEmptySlice,
			slice: []int{},
			want:  []int{},
		}, {
			name:  tcSingleElement,
			slice: []int{1},
			want:  []int{1},
		}, {
			name:  tcMultipleElements,
			slice: []int{1, 2, 3},
			want:  []int{3, 2, 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := Reverse(tc.slice)
			assert.ElementsMatch(t, tc.want, res)
			assert.NotSame(t, &tc.slice, &res)
		})
	}
}

func TestReverseOfStruct(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string
		Age  int
	}

	testCases := []struct {
		name  string
		slice []Person
		want  []Person
	}{
		{
			name:  tcEmptySlice,
			slice: []Person{},
			want:  []Person{},
		}, {
			name:  tcSingleElement,
			slice: []Person{{Name: personJohn, Age: 30}},
			want:  []Person{{Name: personJohn, Age: 30}},
		}, {
			name:  tcMultipleElements,
			slice: []Person{{Name: personJohn, Age: 30}, {Name: "Jane", Age: 25}},
			want:  []Person{{Name: "Jane", Age: 25}, {Name: personJohn, Age: 30}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := Reverse(tc.slice)
			assert.ElementsMatch(t, tc.want, res)
			assert.NotSame(t, &tc.slice, &res)
		})
	}
}

func ExampleReverse() {
	res1 := Reverse([]int{1, 2, 3, 2, 4})
	fmt.Println(res1)

	res2 := Reverse([]string{"a", "b", "c"})
	fmt.Println(res2)
	// Output:
	// [4 2 3 2 1]
	// [c b a]
}
