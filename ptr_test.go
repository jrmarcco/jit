package jit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// personJohn 为 TestPtrOfStruct / TestDePtrOfStruct 复用的样例人名 ( goconst )。
const personJohn = "John"

func TestPtrOfInt(t *testing.T) {
	t.Parallel()

	ptr := Ptr(1)
	assert.Equal(t, 1, *ptr)
}

func TestPtrOfString(t *testing.T) {
	t.Parallel()

	ptr := Ptr("hello")
	assert.Equal(t, "hello", *ptr)
}

func TestPtrOfStruct(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string
		Age  int
	}

	ptr := Ptr(Person{Name: personJohn, Age: 20})
	assert.Equal(t, Person{Name: personJohn, Age: 20}, *ptr)
}

func TestDePtrOfInt(t *testing.T) {
	t.Parallel()

	ptr := Ptr(1)
	assert.Equal(t, 1, DePtr(ptr))
}

func TestDePtrOfString(t *testing.T) {
	t.Parallel()

	ptr := Ptr("hello")
	assert.Equal(t, "hello", DePtr(ptr))
}

func TestDePtrOfStruct(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string
		Age  int
	}

	ptr := Ptr(Person{Name: personJohn, Age: 20})
	assert.Equal(t, Person{Name: personJohn, Age: 20}, DePtr(ptr))
}
