// Package queue provides a FIFO queue backed by a growable ring buffer.
package queue

import (
	"strings"

	"github.com/MorePeanuts/gods/container"
	"github.com/MorePeanuts/gods/deque"
)

// Interface defines the operations of a FIFO queue.
// Operations that return an element use the zero value and false when empty.
type Interface[T any] interface {
	Enqueue(val T)
	Dequeue() (T, bool)
	Peek() (T, bool)

	container.Container[T]
}

var _ Interface[int] = (*Queue[int])(nil)

// Queue is a FIFO queue backed by a growable ring buffer.
// Its zero value is empty and ready to use.
// A Queue must not be copied after first use. Use Clone to make a copy.
type Queue[T any] struct {
	dq deque.Deque[T]
}

// New returns an empty queue.
func New[T any]() *Queue[T] {
	return &Queue[T]{dq: *deque.New[T]()}
}

// NewWithBufSize returns an empty queue with bufSize rounded up to a power of two.
// A size of zero delays allocation until the first insertion. It panics if
// bufSize is negative or cannot be rounded up to a power of two that fits in int.
func NewWithBufSize[T any](bufSize int) *Queue[T] {
	return &Queue[T]{dq: *deque.NewWithBufSize[T](bufSize)}
}

// Enqueue adds val to the back of the queue.
func (qu *Queue[T]) Enqueue(val T) {
	qu.dq.PushBack(val)
}

// Dequeue removes and returns the front element, or the zero value and false if empty.
func (qu *Queue[T]) Dequeue() (T, bool) {
	return qu.dq.PopFront()
}

// Peek returns the next element to be removed, or the zero value and false if empty.
func (qu *Queue[T]) Peek() (T, bool) {
	return qu.dq.Front()
}

// IsEmpty reports whether the queue has no elements.
func (qu *Queue[T]) IsEmpty() bool {
	return qu.dq.IsEmpty()
}

// Len returns the number of elements.
func (qu *Queue[T]) Len() int {
	return qu.dq.Len()
}

// Clear removes all elements and retains the buffer.
func (qu *Queue[T]) Clear() {
	qu.dq.Clear()
}

// String formats the elements from front to back.
func (qu *Queue[T]) String() string {
	strs := qu.dq.String()
	strs = strings.Replace(strs, "deque.Deque", "queue.Queue", 1)
	return strs
}

// Clone returns a shallow copy of qu with independent storage.
// Elements are copied by assignment.
func (qu *Queue[T]) Clone() *Queue[T] {
	return &Queue[T]{dq: *qu.dq.Clone()}
}
