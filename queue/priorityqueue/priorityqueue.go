// Package priorityqueue provides a priority queue backed by a binary heap.
package priorityqueue

import (
	"cmp"
	"strings"

	"github.com/MorePeanuts/gods/queue"
	"github.com/MorePeanuts/gods/tree/binaryheap"
)

var _ queue.Interface[int] = (*Queue[int])(nil)

// Queue removes elements in comparator order. Elements with equal priority
// may be removed in any order. Use New or NewWith before enqueuing elements.
// A Queue must not be copied after first use. Use Clone for an independent copy.
type Queue[T any] struct {
	heap binaryheap.Heap[T]
}

// New builds a priority queue from vals in O(n) time, using cmp.Compare.
// The smallest element is dequeued first.
func New[T cmp.Ordered](vals ...T) *Queue[T] {
	return &Queue[T]{heap: *binaryheap.New(vals...)}
}

// NewWith builds a priority queue from vals in O(n) time. If comparator(a, b)
// is negative, a has higher priority than b. The comparator must define a
// consistent ordering. NewWith panics if comparator is nil.
func NewWith[T any](comparator func(a, b T) int, vals ...T) *Queue[T] {
	return &Queue[T]{heap: *binaryheap.NewWith(comparator, vals...)}
}

// Enqueue adds val to the queue. It panics if the comparator is uninitialized.
func (qu *Queue[T]) Enqueue(val T) {
	qu.heap.Push(val)
}

// Dequeue removes and returns the highest-priority element, or the zero value
// and false if the queue is empty.
func (qu *Queue[T]) Dequeue() (T, bool) {
	return qu.heap.Pop()
}

// Peek returns the highest-priority element without removing it, or the zero
// value and false if the queue is empty.
func (qu *Queue[T]) Peek() (T, bool) {
	return qu.heap.Peek()
}

// IsEmpty reports whether the queue has no elements.
func (qu *Queue[T]) IsEmpty() bool {
	return qu.heap.IsEmpty()
}

// Len returns the number of elements in the queue.
func (qu *Queue[T]) Len() int {
	return qu.heap.Len()
}

// Clear removes all elements and retains the comparator and storage capacity.
func (qu *Queue[T]) Clear() {
	qu.heap.Clear()
}

// String renders the underlying heap as a tree.
func (qu *Queue[T]) String() string {
	return strings.Replace(qu.heap.String(), "binaryheap.Heap", "priorityqueue.Queue", 1)
}

// Clone returns a shallow copy with independent storage and the same
// comparator. Elements are copied by assignment.
func (qu *Queue[T]) Clone() *Queue[T] {
	return &Queue[T]{heap: *qu.heap.Clone()}
}
