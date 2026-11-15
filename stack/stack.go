// Package stack provides a LIFO stack backed by a slice.
package stack

import (
	"fmt"
	"slices"
	"strings"

	"github.com/MorePeanuts/gods/container"
)

// Interface defines the operations of a LIFO stack.
// Operations that return an element use the zero value and false when empty.
type Interface[T any] interface {
	Push(val T)
	Pop() (T, bool)
	Peek() (T, bool)

	container.Container[T]
}

var _ Interface[int] = (*Stack[int])(nil)

// Stack is a LIFO stack backed by a slice.
// Its zero value is empty and ready to use.
// A Stack must not be copied after first use. Use Clone to make a copy.
type Stack[T any] struct {
	elements []T
}

// New returns an empty stack.
func New[T any]() *Stack[T] {
	return &Stack[T]{}
}

// Push adds val to the top of the stack.
func (stk *Stack[T]) Push(val T) {
	stk.elements = append(stk.elements, val)
}

// Pop removes and returns the top element, or the zero value and false if empty.
func (stk *Stack[T]) Pop() (T, bool) {
	if len(stk.elements) == 0 {
		return *new(T), false
	}
	last := len(stk.elements) - 1
	val := stk.elements[last]
	stk.elements[last] = *new(T)
	stk.elements = stk.elements[:last]
	return val, true
}

// Peek returns the next element to be removed, or the zero value and false if empty.
func (stk *Stack[T]) Peek() (T, bool) {
	if len(stk.elements) == 0 {
		return *new(T), false
	}
	return stk.elements[len(stk.elements)-1], true
}

// IsEmpty reports whether the stack has no elements.
func (stk *Stack[T]) IsEmpty() bool {
	return len(stk.elements) == 0
}

// Len returns the number of elements.
func (stk *Stack[T]) Len() int {
	return len(stk.elements)
}

// Clear removes all elements and retains the buffer.
func (stk *Stack[T]) Clear() {
	clear(stk.elements)
	stk.elements = stk.elements[:0]
}

// String formats the elements from bottom to top.
func (stk *Stack[T]) String() string {
	vals := make([]string, 0, len(stk.elements))
	for _, val := range stk.elements {
		vals = append(vals, fmt.Sprintf("%v", val))
	}
	return "stack.Stack{" + strings.Join(vals, ", ") + "}"
}

// Clone returns a shallow copy of stk with independent storage.
// Elements are copied by assignment.
func (stk *Stack[T]) Clone() *Stack[T] {
	return &Stack[T]{elements: slices.Clone(stk.elements)}
}
