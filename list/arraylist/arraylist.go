// Package arraylist provides a generic list backed by a resizable slice.
package arraylist

import (
	"fmt"
	"slices"
	"strings"

	"github.com/MorePeanuts/gods/list"
)

var _ list.Interface[int] = (*List[int])(nil)

// List is a mutable, zero-indexed list backed by a resizable slice.
// Its zero value is an empty list ready to use.
// A List must not be copied after first use. Use Clone to copy its elements.
type List[T any] struct {
	elements []T
}

// New returns a list initialized with vals in argument order.
func New[T any](vals ...T) *List[T] {
	lst := &List[T]{}
	if len(vals) > 0 {
		lst.Append(vals...)
	}
	return lst
}

// Get returns the element at i. It panics if i is outside [0, Len()).
func (lst *List[T]) Get(i int) T {
	if !lst.withinRange(i) {
		panic(fmt.Sprintf("arraylist.List.Get: index %d out of range for length %d", i, len(lst.elements)))
	}
	return lst.elements[i]
}

// Set replaces the element at i. It panics if i is outside [0, Len()).
func (lst *List[T]) Set(i int, val T) {
	if !lst.withinRange(i) {
		panic(fmt.Sprintf("arraylist.List.Set: index %d out of range for length %d", i, len(lst.elements)))
	}
	lst.elements[i] = val
}

// Append adds vals to the end of the list in argument order.
func (lst *List[T]) Append(vals ...T) {
	lst.elements = append(lst.elements, vals...)
}

// RemoveAt deletes and returns the element at i. It panics if i is outside [0, Len()).
func (lst *List[T]) RemoveAt(i int) T {
	if !lst.withinRange(i) {
		panic(fmt.Sprintf("arraylist.List.RemoveAt: index %d out of range for length %d", i, len(lst.elements)))
	}
	element := lst.elements[i]
	lst.elements = slices.Delete(lst.elements, i, i+1)
	return element
}

// SortBy sorts the list according to compare, which must define a strict weak ordering.
func (lst *List[T]) SortBy(compare func(a, b T) int) {
	if len(lst.elements) < 2 {
		return
	}
	slices.SortFunc(lst.elements, compare)
}

// SortStableBy sorts like SortBy and preserves the original order of equal elements.
func (lst *List[T]) SortStableBy(compare func(a, b T) int) {
	if len(lst.elements) < 2 {
		return
	}
	slices.SortStableFunc(lst.elements, compare)
}

// Swap exchanges the elements at i and j. It panics if either index is outside [0, Len()).
func (lst *List[T]) Swap(i, j int) {
	if !lst.withinRange(i) {
		panic(fmt.Sprintf("arraylist.List.Swap: index %d out of range for length %d", i, len(lst.elements)))
	}
	if !lst.withinRange(j) {
		panic(fmt.Sprintf("arraylist.List.Swap: index %d out of range for length %d", j, len(lst.elements)))
	}
	lst.elements[i], lst.elements[j] = lst.elements[j], lst.elements[i]
}

// Insert inserts vals in argument order before i, appending if i equals Len().
// It panics if i is outside [0, Len()].
func (lst *List[T]) Insert(i int, vals ...T) {
	if !lst.withinRange(i) && i != len(lst.elements) {
		panic(fmt.Sprintf("arraylist.List.Insert: index %d out of range [0, %d]", i, len(lst.elements)))
	}
	lst.elements = slices.Insert(lst.elements, i, vals...)
}

// IndexFunc returns the index of the first element for which match returns
// true, or -1 if none match.
func (lst *List[T]) IndexFunc(match func(T) bool) int {
	return slices.IndexFunc(lst.elements, match)
}

// Clone returns a shallow copy of lst with independent storage.
// Elements are copied by assignment.
func (lst *List[T]) Clone() *List[T] {
	return &List[T]{elements: slices.Clone(lst.elements)}
}

// IsEmpty reports whether the list has no elements.
func (lst *List[T]) IsEmpty() bool {
	return len(lst.elements) == 0
}

// Len returns the number of elements.
func (lst *List[T]) Len() int {
	return len(lst.elements)
}

// Clear removes all elements and retains the list capacity.
func (lst *List[T]) Clear() {
	clear(lst.elements)
	lst.elements = lst.elements[:0]
}

// String returns a string representation of the list in element order.
func (lst *List[T]) String() string {
	vals := make([]string, 0, len(lst.elements))
	for _, val := range lst.elements {
		vals = append(vals, fmt.Sprintf("%v", val))
	}
	return "arraylist.List{" + strings.Join(vals, ", ") + "}"
}

// ShrinkToFit reduces the capacity to the current length.
func (lst *List[T]) ShrinkToFit() {
	lst.ShrinkTo(len(lst.elements))
}

// ShrinkTo reduces the capacity to the greater of n and the current length when that value is below the current capacity.
func (lst *List[T]) ShrinkTo(n int) {
	target := max(len(lst.elements), n)
	if target < cap(lst.elements) {
		lst.resize(len(lst.elements), target)
	}
}

func (lst *List[T]) resize(length, cap int) {
	newElements := make([]T, length, cap)
	copy(newElements, lst.elements)
	lst.elements = newElements
}

func (lst *List[T]) withinRange(i int) bool {
	return i >= 0 && i < len(lst.elements)
}
