// Package hashmap provides a hash-based map.
package hashmap

import (
	"fmt"
	stdmaps "maps"
	"strings"

	"github.com/MorePeanuts/gods/maps"
)

// Map maps keys to values. Its zero value is an empty map ready to use.
// Copying a Map by assignment can share its storage. Use Clone for an independent copy.
type Map[K comparable, V any] struct {
	m map[K]V
}

var _ maps.Interface[string, int] = (*Map[string, int])(nil)

// New returns an empty map.
func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{make(map[K]V)}
}

// Get returns the value associated with key and whether key is present.
// If key is absent, Get returns the zero value and false.
func (mp *Map[K, V]) Get(key K) (V, bool) {
	val, found := mp.m[key]
	return val, found
}

// Put associates val with key and returns the previous value and true if key
// was already present. For a new key, it returns the zero value and false.
func (mp *Map[K, V]) Put(key K, val V) (old V, replaced bool) {
	if mp.m == nil {
		mp.m = make(map[K]V)
	}
	old, replaced = mp.m[key]
	mp.m[key] = val
	return
}

// Remove deletes key and returns its previous value and true.
// If key is absent, Remove returns the zero value and false.
func (mp *Map[K, V]) Remove(key K) (old V, removed bool) {
	old, removed = mp.m[key]
	delete(mp.m, key)
	return
}

// Keys returns an independent slice of the keys in unspecified order.
func (mp *Map[K, V]) Keys() []K {
	keys := make([]K, 0, len(mp.m))
	for k := range mp.m {
		keys = append(keys, k)
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
}

// String returns a string representation of the map. Entry order is unspecified.
func (mp *Map[K, V]) String() string {
	lines := make([]string, 0, len(mp.m)+2)
	lines = append(lines, "hashmap.Map{")
	for k, v := range mp.m {
		lines = append(lines, fmt.Sprintf("\t%v: %v", k, v))
	}
	lines = append(lines, "}")
	return strings.Join(lines, "\n")
}

// Clone returns a shallow copy of mp with independent storage.
// Keys and values are copied by assignment.
func (mp *Map[K, V]) Clone() *Map[K, V] {
	return &Map[K, V]{stdmaps.Clone(mp.m)}
}
