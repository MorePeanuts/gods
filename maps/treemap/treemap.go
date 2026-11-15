// Package treemap provides a map ordered by its keys.
package treemap

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/MorePeanuts/gods/maps"
	"github.com/MorePeanuts/gods/tree/rbtree"
)

var _ maps.Interface[int, string] = (*Map[int, string])(nil)

// Map stores key-value pairs in a left-leaning red-black tree.
// Keys that compare equal represent the same entry. Use New or
// NewWithComparator before inserting into a Map. A Map must not be copied
// after first use. Use Clone for an independent copy.
type Map[K comparable, V any] struct {
	tree rbtree.Tree[K, V]
}

// New returns an empty map ordered by cmp.Compare.
func New[K cmp.Ordered, V any]() *Map[K, V] {
	return NewWithComparator[K, V](cmp.Compare[K])
}

// NewWithComparator returns an empty map ordered by compare. The comparator
// must define a consistent ordering and return zero for keys treated as equal.
// NewWithComparator panics if compare is nil.
func NewWithComparator[K comparable, V any](compare func(a, b K) int) *Map[K, V] {
	return &Map[K, V]{tree: *rbtree.NewWithComparator[K, V](compare)}
}

// Get returns the value associated with key and whether key is present.
// If key is absent, it returns the zero value and false.
func (mp *Map[K, V]) Get(key K) (V, bool) {
	return mp.tree.Get(key)
}

// Put associates val with key. It returns the previous value and true if key
// was present, or the zero value and false for a new key. An existing entry
// retains its original key. Put panics if no comparator has been initialized.
func (mp *Map[K, V]) Put(key K, val V) (old V, replaced bool) {
	return mp.tree.Put(key, val)
}

// Remove deletes key and returns its previous value and true.
// If key is absent, it returns the zero value and false.
func (mp *Map[K, V]) Remove(key K) (old V, removed bool) {
	return mp.tree.Remove(key)
}

// Keys returns an independent slice of keys in comparator order.
func (mp *Map[K, V]) Keys() []K {
	keys := make([]K, 0, mp.tree.Len())
	mp.tree.Ascend(func(key K, _ V) bool {
		keys = append(keys, key)
		return true
	})
	return keys
}

// First returns the first entry in comparator order, or zero values and false
// if the map is empty.
func (mp *Map[K, V]) First() (K, V, bool) {
	return mp.tree.First()
}

// Last returns the last entry in comparator order, or zero values and false
// if the map is empty.
func (mp *Map[K, V]) Last() (K, V, bool) {
	return mp.tree.Last()
}

// Floor returns the last entry whose key is at or before key in comparator
// order, or zero values and false if no such entry exists.
func (mp *Map[K, V]) Floor(key K) (K, V, bool) {
	return mp.tree.Floor(key)
}

// Ceiling returns the first entry whose key is at or after key in comparator
// order, or zero values and false if no such entry exists.
func (mp *Map[K, V]) Ceiling(key K) (K, V, bool) {
	return mp.tree.Ceiling(key)
}

// Lower returns the last entry whose key is strictly before key in comparator
// order, or zero values and false if no such entry exists.
func (mp *Map[K, V]) Lower(key K) (K, V, bool) {
	return mp.tree.Lower(key)
}

// Higher returns the first entry whose key is strictly after key in comparator
// order, or zero values and false if no such entry exists.
func (mp *Map[K, V]) Higher(key K) (K, V, bool) {
	return mp.tree.Higher(key)
}

// KeysInRange returns an independent slice of keys in the half-open interval
// [from, to), in comparator order. The result is empty if from compares
// equal to or after to. The traversal takes O(log n + k) time for k results.
func (mp *Map[K, V]) KeysInRange(from, to K) []K {
	keys := make([]K, 0)
	mp.tree.AscendRange(from, to, func(key K, _ V) bool {
		keys = append(keys, key)
		return true
	})
	return keys
}

// IsEmpty reports whether the map has no entries.
func (mp *Map[K, V]) IsEmpty() bool {
	return mp.tree.IsEmpty()
}

// Len returns the number of entries in the map.
func (mp *Map[K, V]) Len() int {
	return mp.tree.Len()
}

// Clear removes all entries and retains the comparator.
func (mp *Map[K, V]) Clear() {
	mp.tree.Clear()
}

// String returns a string representation of the map in comparator order.
func (mp *Map[K, V]) String() string {
	lines := make([]string, 0, mp.tree.Len()+2)
	lines = append(lines, "treemap.Map{")
	mp.tree.Ascend(func(key K, val V) bool {
		lines = append(lines, fmt.Sprintf("\t%v: %v", key, val))
		return true
	})
	lines = append(lines, "}")
	return strings.Join(lines, "\n")
}

// Clone returns a shallow copy with independent nodes and the same comparator.
// Keys and values are copied by assignment.
func (mp *Map[K, V]) Clone() *Map[K, V] {
	return &Map[K, V]{tree: *mp.tree.Clone()}
}
