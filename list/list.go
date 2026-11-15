package list

import "github.com/MorePeanuts/gods/container"

type List[T comparable] interface {
	Get(i int) T
	Set(i int, val T)
	Append(vals ...T)
	Remove(i int) T
	Contains(val T) bool
	ContainsAll(vals ...T) bool
	SortBy(compare func(a, b T) int)
	SortStableBy(compare func(a, b T) int)
	Swap(i, j int)
	Insert(i int, vals ...T)
	IndexOf(val T) int

	container.Container[T]
	// IsEmpty() bool
	// Len() int
	// Clear()
	// String() string
}
