// Package linkedlist provides a generic list backed by doubly linked nodes.
package linkedlist

import (
	"fmt"
	"slices"
	"strings"

	"github.com/MorePeanuts/gods/deque"
	"github.com/MorePeanuts/gods/list"
	"github.com/MorePeanuts/gods/queue"
	"github.com/MorePeanuts/gods/stack"
)

var (
	_ list.Interface[int]  = (*List[int])(nil)
	_ deque.Interface[int] = (*List[int])(nil)
	_ queue.Interface[int] = (*List[int])(nil)
	_ stack.Interface[int] = (*List[int])(nil)
)

// List is a mutable, zero-indexed list backed by doubly linked nodes.
// Its zero value is an empty list ready to use.
// A List must not be copied after first use. Use Clone to copy its elements.
type List[T any] struct {
	head   *node[T]
	length int
}

type node[T any] struct {
	val  T
	prev *node[T]
	next *node[T]
}

// New returns a list initialized with vals in argument order.
func New[T any](vals ...T) *List[T] {
	lst := &List[T]{}
	lst.Append(vals...)
	return lst
}

// initHead allocates the sentinel before the first insertion.
func (lst *List[T]) initHead() {
	if lst.head == nil {
		head := &node[T]{}
		head.next = head
		head.prev = head
		lst.head = head
	}
}

// Get returns the element at i. It panics if i is outside [0, Len()).
func (lst *List[T]) Get(i int) T {
	if !lst.withinRange(i) {
		panic(fmt.Sprintf("linkedlist.List.Get: index %d out of range for length %d", i, lst.length))
	}
	p := lst.head
	if lst.length-i > i+1 {
		for range i + 1 {
			p = p.next
		}
	} else {
		for range lst.length - i {
			p = p.prev
		}
	}
	return p.val
}

// Set replaces the element at i. It panics if i is outside [0, Len()).
func (lst *List[T]) Set(i int, val T) {
	if !lst.withinRange(i) {
		panic(fmt.Sprintf("linkedlist.List.Set: index %d out of range for length %d", i, lst.length))
	}
	p := lst.head
	if lst.length-i > i+1 {
		for range i + 1 {
			p = p.next
		}
	} else {
		for range lst.length - i {
			p = p.prev
		}
	}
	p.val = val
}

// Append adds vals to the end of the list in argument order.
func (lst *List[T]) Append(vals ...T) {
	if len(vals) == 0 {
		return
	}
	lst.initHead()
	for _, val := range vals {
		newNode := &node[T]{val, lst.head.prev, lst.head}
		lst.head.prev.next = newNode
		lst.head.prev = newNode
		lst.length++
	}
}

// RemoveAt deletes and returns the element at i. It panics if i is outside [0, Len()).
func (lst *List[T]) RemoveAt(i int) T {
	if !lst.withinRange(i) {
		panic(fmt.Sprintf("linkedlist.List.RemoveAt: index %d out of range for length %d", i, lst.length))
	}
	p := lst.head
	if lst.length-i > i+1 {
		for range i + 1 {
			p = p.next
		}
	} else {
		for range lst.length - i {
			p = p.prev
		}
	}
	p.next.prev = p.prev
	p.prev.next = p.next
	lst.length--
	return p.val
}

// SortBy sorts the list according to compare, which must define a strict weak ordering.
func (lst *List[T]) SortBy(compare func(a, b T) int) {
	if lst.length < 2 {
		return
	}
	vals := make([]T, 0, lst.length)
	for p := lst.head.next; p != lst.head; p = p.next {
		vals = append(vals, p.val)
	}
	slices.SortFunc(vals, compare)
	lst.Clear()
	lst.Append(vals...)
}

// SortStableBy sorts like SortBy and preserves the original order of equal elements.
func (lst *List[T]) SortStableBy(compare func(a, b T) int) {
	if lst.length < 2 {
		return
	}
	vals := make([]T, 0, lst.length)
	for p := lst.head.next; p != lst.head; p = p.next {
		vals = append(vals, p.val)
	}
	slices.SortStableFunc(vals, compare)
	lst.Clear()
	lst.Append(vals...)
}

// Swap exchanges the elements at i and j. It panics if either index is outside [0, Len()).
func (lst *List[T]) Swap(i, j int) {
	if !lst.withinRange(i) {
		panic(fmt.Sprintf("linkedlist.List.Swap: index %d out of range for length %d", i, lst.length))
	}
	if !lst.withinRange(j) {
		panic(fmt.Sprintf("linkedlist.List.Swap: index %d out of range for length %d", j, lst.length))
	}
	if i == j {
		return
	}
	var pi, pj *node[T]
	for k, p := 0, lst.head.next; pi == nil || pj == nil; k, p = k+1, p.next {
		switch k {
		case i:
			pi = p
		case j:
			pj = p
		}
	}
	pi.val, pj.val = pj.val, pi.val
}

// Insert inserts vals in argument order before i, appending if i equals Len().
// It panics if i is outside [0, Len()].
func (lst *List[T]) Insert(i int, vals ...T) {
	if !lst.withinRange(i) && i != lst.length {
		panic(fmt.Sprintf("linkedlist.List.Insert: index %d out of range [0, %d]", i, lst.length))
	}
	if len(vals) == 0 {
		return
	}
	lst.initHead()
	p := lst.head
	if lst.length-i+1 > i {
		for range i {
			p = p.next
		}
	} else {
		for range lst.length - i + 1 {
			p = p.prev
		}
	}
	for _, val := range vals {
		newNode := &node[T]{val, p, p.next}
		newNode.prev.next = newNode
		newNode.next.prev = newNode
		p = newNode
		lst.length++
	}
}

// IndexFunc returns the index of the first element for which match returns
// true, or -1 if none match.
func (lst *List[T]) IndexFunc(match func(T) bool) int {
	if lst.head == nil {
		return -1
	}
	for i, p := 0, lst.head.next; p != lst.head; i, p = i+1, p.next {
		if match(p.val) {
			return i
		}
	}
	return -1
}

// Clone returns a shallow copy of lst with independent storage.
// Elements are copied by assignment.
func (lst *List[T]) Clone() *List[T] {
	cloned := &List[T]{}
	if lst.head != nil {
		for p := lst.head.next; p != lst.head; p = p.next {
			cloned.Append(p.val)
		}
	}
	return cloned
}

// IsEmpty reports whether the list has no elements.
func (lst *List[T]) IsEmpty() bool {
	return lst.length == 0
}

// Len returns the number of elements.
func (lst *List[T]) Len() int {
	return lst.length
}

// Clear removes all elements.
func (lst *List[T]) Clear() {
	if lst.head == nil {
		return
	}
	lst.head.next = lst.head
	lst.head.prev = lst.head
	lst.length = 0
}

// String returns a string representation of the list in element order.
func (lst *List[T]) String() string {
	if lst.head == nil {
		return "linkedlist.List{}"
	}
	vals := make([]string, 0, lst.length)
	for n := lst.head.next; n != lst.head; n = n.next {
		vals = append(vals, fmt.Sprintf("%v", n.val))
	}
	return "linkedlist.List{" + strings.Join(vals, ", ") + "}"
}

// Prepend adds vals to the front of the list in argument order.
func (lst *List[T]) Prepend(vals ...T) {
	lst.Insert(0, vals...)
}

// PushFront inserts val at the head of the list.
func (lst *List[T]) PushFront(val T) {
	lst.Prepend(val)
}

// PushBack inserts val at the tail of the list.
func (lst *List[T]) PushBack(val T) {
	lst.Append(val)
}

// Front returns the head element, or the zero value and false if empty.
func (lst *List[T]) Front() (T, bool) {
	if lst.IsEmpty() {
		return *new(T), false
	}
	return lst.head.next.val, true
}

// Back returns the tail element, or the zero value and false if empty.
func (lst *List[T]) Back() (T, bool) {
	if lst.IsEmpty() {
		return *new(T), false
	}
	return lst.head.prev.val, true
}

// PopFront removes and returns the head element, or the zero value and false if empty.
func (lst *List[T]) PopFront() (T, bool) {
	if lst.IsEmpty() {
		return *new(T), false
	}
	return lst.RemoveAt(0), true
}

// PopBack removes and returns the tail element, or the zero value and false if empty.
func (lst *List[T]) PopBack() (T, bool) {
	if lst.IsEmpty() {
		return *new(T), false
	}
	return lst.RemoveAt(lst.length - 1), true
}

// Enqueue inserts val at the tail of the list.
func (lst *List[T]) Enqueue(val T) {
	lst.PushBack(val)
}

// Dequeue removes and returns the head element, or the zero value and false if empty.
func (lst *List[T]) Dequeue() (T, bool) {
	return lst.PopFront()
}

// Push places val at the top of the stack, which is the head of the list.
func (lst *List[T]) Push(val T) {
	lst.PushFront(val)
}

// Pop removes and returns the stack top, or the zero value and false if empty.
func (lst *List[T]) Pop() (T, bool) {
	return lst.PopFront()
}

// Peek returns the head element, which is both the queue front and stack top.
// It returns the zero value and false if empty.
func (lst *List[T]) Peek() (T, bool) {
	return lst.Front()
}

func (lst *List[T]) withinRange(i int) bool {
	return i >= 0 && i < lst.length
}
