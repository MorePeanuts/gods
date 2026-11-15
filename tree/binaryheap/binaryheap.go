// Package binaryheap provides a slice-backed binary heap.
package binaryheap

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/MorePeanuts/gods/tree"
)

// Heap is a complete binary tree whose root holds the element that comes first
// according to its comparator. Use New or NewWith before pushing elements.
// A Heap must not be copied after first use. Use Clone for an independent copy.
type Heap[T any] struct {
	array      []T
	comparator func(a, b T) int
}

var _ tree.Heap[int] = (*Heap[int])(nil)

// New builds a min-heap from vals using cmp.Compare in O(n) time.
func New[T cmp.Ordered](vals ...T) *Heap[T] {
	hp := &Heap[T]{make([]T, 0, len(vals)), cmp.Compare[T]}
	hp.array = append(hp.array, vals...)
	for i := (len(hp.array) - 2) >> 1; i >= 0; i-- {
		hp.bubbleDownAt(i)
	}
	return hp
}

// NewWith builds a heap from vals in O(n) time. If comparator(a, b) is
// negative, a has higher priority than b. The comparator must define a
// consistent ordering. Elements with equal priority may be popped in any order.
// NewWith panics if comparator is nil.
func NewWith[T any](comparator func(a, b T) int, vals ...T) *Heap[T] {
	if comparator == nil {
		panic("binaryheap.Heap: nil comparator")
	}
	hp := &Heap[T]{make([]T, 0, len(vals)), comparator}
	hp.array = append(hp.array, vals...)
	for i := (len(hp.array) - 2) >> 1; i >= 0; i-- {
		hp.bubbleDownAt(i)
	}
	return hp
}

// Push adds val to the heap. It panics if the comparator is uninitialized.
func (hp *Heap[T]) Push(val T) {
	if hp.comparator == nil {
		panic("binaryheap.Heap: initialize with New or NewWith before pushing")
	}
	hp.array = append(hp.array, val)
	hp.bubbleUp()
}

// Pop removes and returns the highest-priority element, or the zero value
// and false if the heap is empty.
func (hp *Heap[T]) Pop() (val T, ok bool) {
	if len(hp.array) == 0 {
		return *new(T), false
	}
	val, ok = hp.array[0], true
	hp.array[0] = hp.array[len(hp.array)-1]
	hp.array[len(hp.array)-1] = *new(T)
	hp.array = hp.array[:len(hp.array)-1]
	if len(hp.array) != 0 {
		hp.bubbleDown()
	}
	return
}

// Peek returns the highest-priority element without removing it, or the zero
// value and false if the heap is empty.
func (hp *Heap[T]) Peek() (T, bool) {
	if len(hp.array) == 0 {
		return *new(T), false
	}
	return hp.array[0], true
}

// IsEmpty reports whether the heap has no elements.
func (hp *Heap[T]) IsEmpty() bool {
	return len(hp.array) == 0
}

// Len returns the number of elements in the heap.
func (hp *Heap[T]) Len() int {
	return len(hp.array)
}

// Clear removes all elements and retains the comparator and storage capacity.
func (hp *Heap[T]) Clear() {
	clear(hp.array)
	hp.array = hp.array[:0]
}

// String renders the complete binary tree rooted at the highest-priority
// element, with each left child shown before its right sibling.
func (hp *Heap[T]) String() string {
	if hp.IsEmpty() {
		return "binaryheap.Heap{}"
	}

	var out strings.Builder
	fmt.Fprintf(&out, "binaryheap.Heap{\n  %v\n", hp.array[0])
	var writeChildren func(parent int, prefix string)
	writeChildren = func(parent int, prefix string) {
		if parent >= len(hp.array)/2 {
			return
		}
		left := 2*parent + 1
		right := left + 1
		if right == len(hp.array) {
			fmt.Fprintf(&out, "%s└── %v\n", prefix, hp.array[left])
			return
		}
		fmt.Fprintf(&out, "%s├── %v\n", prefix, hp.array[left])
		writeChildren(left, prefix+"│   ")
		fmt.Fprintf(&out, "%s└── %v\n", prefix, hp.array[right])
		writeChildren(right, prefix+"    ")
	}
	writeChildren(0, "  ")
	out.WriteByte('}')
	return out.String()
}

// Clone returns a shallow copy with independent storage and the same
// comparator. Elements are copied by assignment.
func (hp *Heap[T]) Clone() *Heap[T] {
	return &Heap[T]{array: slices.Clone(hp.array), comparator: hp.comparator}
}

func (hp *Heap[T]) bubbleDownAt(i int) {
	x := hp.array[i]
	for l := i<<1 + 1; l < len(hp.array); l = i<<1 + 1 {
		r := l + 1
		small := l
		if r < len(hp.array) && hp.comparator(hp.array[r], hp.array[l]) < 0 {
			small = r
		}
		if hp.comparator(hp.array[small], x) >= 0 {
			break
		}
		hp.array[i] = hp.array[small]
		i = small
	}
	hp.array[i] = x
}

func (hp *Heap[T]) bubbleDown() {
	hp.bubbleDownAt(0)
}

func (hp *Heap[T]) bubbleUp() {
	i := len(hp.array) - 1
	x := hp.array[i]
	for p := (i - 1) >> 1; i > 0; p = (i - 1) >> 1 {
		if hp.comparator(hp.array[p], x) <= 0 {
			break
		}
		hp.array[i] = hp.array[p]
		i = p
	}
	hp.array[i] = x
}
