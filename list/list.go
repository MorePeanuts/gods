// Package list defines the common list interface along with search and sorting helpers.
package list

import (
	"cmp"

	"github.com/MorePeanuts/gods/container"
)

// Interface defines the operations of a mutable, ordered list with zero-based indexing.
type Interface[T any] interface {
	// Get returns the element at i.
	Get(i int) T
	// Set replaces the element at i with val.
	Set(i int, val T)
	// Append adds vals to the end in argument order.
	Append(vals ...T)
	// RemoveAt deletes and returns the element at i. It panics if i is outside [0, Len()).
	RemoveAt(i int) T
	// IndexFunc calls match in element order and returns the index of the
	// first element for which match returns true, or -1 if none match.
	IndexFunc(match func(T) bool) int
	// SortBy sorts according to compare.
	SortBy(compare func(a, b T) int)
	// SortStableBy sorts according to compare, preserving the order of equal elements.
	SortStableBy(compare func(a, b T) int)
	// Swap exchanges the elements at i and j.
	Swap(i, j int)
	// Insert inserts vals in argument order before i, or appends them if i equals Len().
	Insert(i int, vals ...T)

	container.Container[T]
}

// IndexOf returns the index of the first element equal to val, or -1 if absent.
func IndexOf[T comparable](lst Interface[T], val T) int {
	return lst.IndexFunc(func(element T) bool { return element == val })
}

// Contains reports whether lst contains an element equal to val.
func Contains[T comparable](lst Interface[T], val T) bool {
	return IndexOf(lst, val) >= 0
}

// ContainsAll reports whether lst contains each of vals at least once.
// Duplicate queries do not require duplicate elements. An empty query returns true.
func ContainsAll[T comparable](lst Interface[T], vals ...T) bool {
	if len(vals) == 0 {
		return true
	}
	if lst.IsEmpty() {
		return false
	}

	remaining := make(map[T]struct{}, len(vals))
	for _, val := range vals {
		remaining[val] = struct{}{}
	}
	return lst.IndexFunc(func(element T) bool {
		delete(remaining, element)
		return len(remaining) == 0
	}) >= 0
}

// Sort sorts lst in ascending order.
func Sort[T cmp.Ordered](lst Interface[T]) {
	lst.SortBy(cmp.Compare[T])
}

// SortStable sorts lst in ascending order, preserving the original order of equal elements.
func SortStable[T cmp.Ordered](lst Interface[T]) {
	lst.SortStableBy(cmp.Compare[T])
}
