// Package deque provides a double-ended queue backed by a growable ring buffer.
package deque

import (
	"fmt"
	"math/bits"
	"strings"

	"github.com/MorePeanuts/gods/container"
)

// Interface defines the operations of a double-ended queue.
type Interface[T any] interface {
	// PushFront inserts val at the front of the deque.
	PushFront(val T)
	// PushBack inserts val at the back of the deque.
	PushBack(val T)
	// PopFront removes and returns the front element, or the zero value and false if empty.
	PopFront() (T, bool)
	// PopBack removes and returns the back element, or the zero value and false if empty.
	PopBack() (T, bool)
	// Front returns the front element without removing it, or the zero value and false if empty.
	Front() (T, bool)
	// Back returns the back element without removing it, or the zero value and false if empty.
	Back() (T, bool)

	container.Container[T]
}

var _ Interface[int] = (*Deque[int])(nil)

// Deque is a double-ended queue backed by a growable ring buffer.
// Its zero value is empty and ready to use.
// A Deque must not be copied after first use. Use Clone to make a copy.
type Deque[T any] struct {
	elements []T
	head     int
	length   int
}

const (
	defaultBufSize = 16
	maxBufSize     = 1 << (bits.UintSize - 2)
)

// New returns an empty deque with the default buffer size.
func New[T any]() *Deque[T] {
	return NewWithBufSize[T](defaultBufSize)
}

// NewWithBufSize returns an empty deque with bufSize rounded up to a power of two.
// A size of zero delays allocation until the first insertion. It panics if
// bufSize is negative or cannot be rounded up to a power of two that fits in int.
func NewWithBufSize[T any](bufSize int) *Deque[T] {
	return &Deque[T]{elements: make([]T, roundedBufSize(bufSize))}
}

func roundedBufSize(size int) int {
	if size < 0 || size > maxBufSize {
		panic("deque: buffer size out of range")
	}
	if size == 0 {
		return 0
	}
	return 1 << bits.Len(uint(size-1))
}

// PushFront inserts val at the front of the deque.
func (dq *Deque[T]) PushFront(val T) {
	dq.growIfFull()
	dq.head = (dq.head - 1) & (len(dq.elements) - 1)
	dq.elements[dq.head] = val
	dq.length++
}

// PushBack inserts val at the back of the deque.
func (dq *Deque[T]) PushBack(val T) {
	dq.growIfFull()
	dq.elements[dq.index(dq.length)] = val
	dq.length++
}

// PopFront removes and returns the front element, or the zero value and false if empty.
func (dq *Deque[T]) PopFront() (T, bool) {
	if dq.length == 0 {
		return *new(T), false
	}
	val := dq.elements[dq.head]
	dq.elements[dq.head] = *new(T)
	dq.head = dq.index(1)
	dq.length--
	return val, true
}

// PopBack removes and returns the back element, or the zero value and false if empty.
func (dq *Deque[T]) PopBack() (T, bool) {
	if dq.length == 0 {
		return *new(T), false
	}
	dq.length--
	idx := dq.index(dq.length)
	val := dq.elements[idx]
	dq.elements[idx] = *new(T)
	return val, true
}

// Front returns the front element, or the zero value and false if empty.
func (dq *Deque[T]) Front() (T, bool) {
	if dq.length == 0 {
		return *new(T), false
	}
	return dq.elements[dq.head], true
}

// Back returns the back element, or the zero value and false if empty.
func (dq *Deque[T]) Back() (T, bool) {
	if dq.length == 0 {
		return *new(T), false
	}
	return dq.elements[dq.index(dq.length-1)], true
}

// IsEmpty reports whether the deque has no elements.
func (dq *Deque[T]) IsEmpty() bool {
	return dq.length == 0
}

// Len returns the number of elements.
func (dq *Deque[T]) Len() int {
	return dq.length
}

// Clear removes all elements and retains the buffer.
func (dq *Deque[T]) Clear() {
	clear(dq.elements)
	dq.head = 0
	dq.length = 0
}

// String formats the elements from front to back.
func (dq *Deque[T]) String() string {
	vals := make([]string, 0, dq.length)
	for i := range dq.length {
		vals = append(vals, fmt.Sprintf("%v", dq.elements[dq.index(i)]))
	}
	return "deque.Deque{" + strings.Join(vals, ", ") + "}"
}

// Clone returns a shallow copy of dq with independent storage.
// Elements are copied by assignment.
func (dq *Deque[T]) Clone() *Deque[T] {
	clone := NewWithBufSize[T](dq.length)
	dq.copyElements(clone.elements)
	clone.length = dq.length
	return clone
}

func (dq *Deque[T]) index(offset int) int {
	return (dq.head + offset) & (len(dq.elements) - 1)
}

func (dq *Deque[T]) growIfFull() {
	size := len(dq.elements)
	if dq.length < size {
		return
	}
	if size == maxBufSize {
		panic("deque: buffer size overflow")
	}
	newSize := size * 2
	if size == 0 {
		newSize = defaultBufSize
	}
	elements := make([]T, newSize)
	dq.copyElements(elements)
	dq.elements = elements
	dq.head = 0
}

// copyElements copies the live elements in front-to-back order.
func (dq *Deque[T]) copyElements(dst []T) {
	first := min(dq.length, len(dq.elements)-dq.head)
	copy(dst, dq.elements[dq.head:dq.head+first])
	copy(dst[first:], dq.elements[:dq.length-first])
}
