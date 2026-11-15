// Package tree defines contracts shared by tree data structures.
package tree

import "github.com/MorePeanuts/gods/container"

// SearchTree is a mutable collection of key-value pairs ordered by a comparator.
// Keys that compare equal identify the same entry.
// Methods that return an entry return zero values and false if none matches.
type SearchTree[K any, V any] interface {
	// Get returns the value for key, or the zero value and false if absent.
	Get(key K) (V, bool)
	// Put associates val with key. If an equivalent key exists, Put retains
	// that key and returns its previous value and true. Otherwise it returns
	// the zero value and false.
	Put(key K, val V) (old V, replaced bool)
	// Remove deletes key and returns its previous value and true, or the zero
	// value and false if absent.
	Remove(key K) (old V, removed bool)

	// First returns the first entry in comparator order.
	First() (K, V, bool)
	// Last returns the last entry in comparator order.
	Last() (K, V, bool)
	// Floor returns the last entry at or before key in comparator order.
	Floor(key K) (K, V, bool)
	// Ceiling returns the first entry at or after key in comparator order.
	Ceiling(key K) (K, V, bool)
	// Lower returns the last entry strictly before key in comparator order.
	Lower(key K) (K, V, bool)
	// Higher returns the first entry strictly after key in comparator order.
	Higher(key K) (K, V, bool)

	// Ascend calls yield in comparator order until it returns false.
	// The callback must not modify the tree.
	Ascend(yield func(K, V) bool)
	// AscendRange calls yield in comparator order for entries in [from, to)
	// until it returns false. It yields nothing if from is at or after to.
	// The callback must not modify the tree.
	AscendRange(from, to K, yield func(K, V) bool)

	container.Container[V]
}

// Heap stores elements in priority order.
// Each implementation defines which element has the highest priority.
type Heap[T any] interface {
	// Push adds val to the heap.
	Push(val T)

	// Peek returns the highest-priority element without removing it,
	// or the zero value and false if the heap is empty.
	Peek() (T, bool)

	// Pop removes and returns the highest-priority element,
	// or the zero value and false if the heap is empty.
	Pop() (T, bool)

	container.Container[T]
}
