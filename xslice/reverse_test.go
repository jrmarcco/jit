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
			name:  "empty slice",
			slice: []int{},
			want:  []int{},
		}, {
			name:  "single element",
			slice: []int{1},
			want:  []int{1},
		}, {
			name:  "multiple elements",
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
			name:  "empty slice",
			slice: []Person{},
			want:  []Person{},
		}, {
			name:  "single element",
			slice: []Person{{Name: "John", Age: 30}},
			want:  []Person{{Name: "John", Age: 30}},
		}, {
			name:  "multiple elements",
			slice: []Person{{Name: "John", Age: 30}, {Name: "Jane", Age: 25}},
			want:  []Person{{Name: "Jane", Age: 25}, {Name: "John", Age: 30}},
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

func TestReverseInPlaceOfInt(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		slice []int
		want  []int
	}{
		{
			name:  "empty slice",
			slice: []int{},
			want:  []int{},
		}, {
			name:  "single element",
			slice: []int{1},
			want:  []int{1},
		}, {
			name:  "multiple elements",
			slice: []int{1, 2, 3},
			want:  []int{3, 2, 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ReverseInPlace(tc.slice)
			assert.ElementsMatch(t, tc.want, tc.slice)
		})
	}
}

func TestReverseInPlaceOfStruct(t *testing.T) {
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
			name:  "empty slice",
			slice: []Person{},
			want:  []Person{},
		}, {
			name:  "single element",
			slice: []Person{{Name: "John", Age: 30}},
			want:  []Person{{Name: "John", Age: 30}},
		}, {
			name:  "multiple elements",
			slice: []Person{{Name: "John", Age: 30}, {Name: "Jane", Age: 25}},
			want:  []Person{{Name: "Jane", Age: 25}, {Name: "John", Age: 30}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ReverseInPlace(tc.slice)
			assert.ElementsMatch(t, tc.want, tc.slice)
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

func ExampleReverseInPlace() {
	s1 := []int{1, 2, 3, 2, 4}
	ReverseInPlace(s1)
	fmt.Println(s1)

	s2 := []string{"a", "b", "c", "g", "d", "e"}
	ReverseInPlace(s2)
	fmt.Println(s2)

	// Output:
	// [4 2 3 2 1]
	// [e d g c b a]
}
