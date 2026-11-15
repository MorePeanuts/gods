// Package treemap provides a map ordered by its keys.
package treemap

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/MorePeanuts/gods/maps"
)

var _ maps.Interface[int, string] = (*Map[int, string])(nil)

// Map stores key-value pairs in a left-leaning red-black tree.
// Keys that compare equal represent the same entry. Use New or
// NewWithComparator before inserting into a Map. A Map must not be copied
// after first use. Use Clone for an independent copy.
type Map[K comparable, V any] struct {
	root    *node[K, V]
	length  int
	compare func(a, b K) int
}

// New returns an empty map ordered by cmp.Compare.
func New[K cmp.Ordered, V any]() *Map[K, V] {
	return NewWithComparator[K, V](cmp.Compare[K])
}

// NewWithComparator returns an empty map ordered by compare. The comparator
// must define a consistent ordering and return zero for keys treated as equal.
// NewWithComparator panics if compare is nil.
func NewWithComparator[K comparable, V any](compare func(a, b K) int) *Map[K, V] {
	if compare == nil {
		panic("treemap: nil comparator")
	}
	return &Map[K, V]{compare: compare}
}

func (mp *Map[K, V]) find(key K) *node[K, V] {
	for n := mp.root; n != nil; {
		switch order := mp.compare(key, n.key); {
		case order < 0:
			n = n.left
		case order > 0:
			n = n.right
		default:
			return n
		}
	}
	return nil
}

// Get returns the value associated with key and whether key is present.
// If key is absent, it returns the zero value and false.
func (mp *Map[K, V]) Get(key K) (V, bool) {
	if n := mp.find(key); n != nil {
		return n.value, true
	}
	return *new(V), false
}

// Put associates val with key. It returns the previous value and true if key
// was present, or the zero value and false for a new key. An existing entry
// retains its original key. Put panics if no comparator has been initialized.
func (mp *Map[K, V]) Put(key K, val V) (old V, replaced bool) {
	if mp.compare == nil {
		panic("treemap: initialize with New or NewWithComparator before inserting")
	}
	mp.root, old, replaced = mp.insert(mp.root, key, val)
	mp.root.red = false
	if !replaced {
		mp.length++
	}
	return old, replaced
}

// Remove deletes key and returns its previous value and true.
// If key is absent, it returns the zero value and false.
func (mp *Map[K, V]) Remove(key K) (old V, removed bool) {
	old, removed = mp.Get(key)
	if !removed {
		return old, false
	}

	// Start deletion with a red root when both children are black so
	// the descent can move a red link toward the node being removed.
	if !isRed(mp.root.left) && !isRed(mp.root.right) {
		mp.root.red = true
	}
	mp.root = mp.remove(mp.root, key)
	if mp.root != nil {
		mp.root.red = false
	}
	mp.length--
	return old, true
}

// Keys returns an independent slice of keys in comparator order.
func (mp *Map[K, V]) Keys() []K {
	keys := make([]K, 0, mp.length)
	mp.root.visit(func(key K, _ V) {
		keys = append(keys, key)
	})
	return keys
}

// IsEmpty reports whether the map has no entries.
func (mp *Map[K, V]) IsEmpty() bool {
	return mp.length == 0
}

// Len returns the number of entries in the map.
func (mp *Map[K, V]) Len() int {
	return mp.length
}

// Clear removes all entries and retains the comparator.
func (mp *Map[K, V]) Clear() {
	mp.root = nil
	mp.length = 0
}

// String returns a string representation of the map in comparator order.
func (mp *Map[K, V]) String() string {
	lines := make([]string, 0, mp.length+2)
	lines = append(lines, "treemap.Map{")
	mp.root.visit(func(key K, val V) {
		lines = append(lines, fmt.Sprintf("\t%v: %v", key, val))
	})
	lines = append(lines, "}")
	return strings.Join(lines, "\n")
}

// Clone returns a shallow copy with independent nodes and the same comparator.
// Keys and values are copied by assignment.
func (mp *Map[K, V]) Clone() *Map[K, V] {
	return &Map[K, V]{root: mp.root.clone(), length: mp.length, compare: mp.compare}
}
