// Package linkedhashmap provides a hash map that preserves key insertion order.
package linkedhashmap

import (
	"fmt"
	"strings"

	"github.com/MorePeanuts/gods/maps"
)

// Map maps keys to values in insertion order. Its zero value is ready to use.
// Copying a Map by assignment can share its storage. Use Clone for an independent copy.
type Map[K comparable, V any] struct {
	m    map[K]*entry[K, V]
	head *entry[K, V] // Sentinel of the circular list, allocated on first insertion.
}

type entry[K comparable, V any] struct {
	key   K
	value V
	prev  *entry[K, V]
	next  *entry[K, V]
}

var _ maps.Interface[string, int] = (*Map[string, int])(nil)

// New returns an empty map.
func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{}
}

// Get returns the value associated with key and whether key is present.
// If key is absent, Get returns the zero value and false.
func (mp *Map[K, V]) Get(key K) (V, bool) {
	e, found := mp.m[key]
	if !found {
		return *new(V), false
	}
	return e.value, true
}

// Put associates val with key and returns the previous value and true if key
// was already present. An existing key keeps its position. A new key is
// appended to the end, and Put returns the zero value and false.
func (mp *Map[K, V]) Put(key K, val V) (old V, replaced bool) {
	if e, found := mp.m[key]; found {
		old := e.value
		e.value = val
		return old, true
	}
	if mp.m == nil {
		mp.m = make(map[K]*entry[K, V])
		mp.head = &entry[K, V]{}
		mp.head.prev = mp.head
		mp.head.next = mp.head
	}

	e := &entry[K, V]{key: key, value: val, prev: mp.head.prev, next: mp.head}
	mp.m[key] = e
	e.prev.next = e
	e.next.prev = e
	return *new(V), false
}

// Remove deletes key and returns its previous value and true.
// If key is absent, Remove returns the zero value and false.
func (mp *Map[K, V]) Remove(key K) (old V, removed bool) {
	e, found := mp.m[key]
	if !found {
		return *new(V), false
	}
	delete(mp.m, key)

	e.prev.next = e.next
	e.next.prev = e.prev
	return e.value, true
}

// Keys returns an independent slice of keys in insertion order.
func (mp *Map[K, V]) Keys() []K {
	keys := make([]K, 0, len(mp.m))
	if mp.head != nil {
		for e := mp.head.next; e != mp.head; e = e.next {
			keys = append(keys, e.key)
		}
	}
	return keys
}

// IsEmpty reports whether the map has no entries.
func (mp *Map[K, V]) IsEmpty() bool {
	return len(mp.m) == 0
}

// Len returns the number of entries in the map.
func (mp *Map[K, V]) Len() int {
	return len(mp.m)
}

// Clear removes all entries from the map.
func (mp *Map[K, V]) Clear() {
	clear(mp.m)
	if mp.head != nil {
		mp.head.prev = mp.head
		mp.head.next = mp.head
	}
}

// String returns a string representation of the map in insertion order.
func (mp *Map[K, V]) String() string {
	lines := make([]string, 0, len(mp.m)+2)
	lines = append(lines, "linkedhashmap.Map{")
	if mp.head != nil {
		for e := mp.head.next; e != mp.head; e = e.next {
			lines = append(lines, fmt.Sprintf("\t%v: %v", e.key, e.value))
		}
	}
	lines = append(lines, "}")
	return strings.Join(lines, "\n")
}

// Clone returns a shallow copy of mp with independent storage and the same
// insertion order. Keys and values are copied by assignment.
func (mp *Map[K, V]) Clone() *Map[K, V] {
	cloned := New[K, V]()
	if mp.head != nil {
		for e := mp.head.next; e != mp.head; e = e.next {
			cloned.Put(e.key, e.value)
		}
	}
	return cloned
}
