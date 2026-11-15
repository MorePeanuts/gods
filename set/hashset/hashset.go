// Package hashset provides a hash-based set.
package hashset

import (
	"fmt"
	"maps"
	"strings"

	"github.com/MorePeanuts/gods/set"
)

var _ set.Interface[int] = (*Set[int])(nil)

// Set is an unordered collection of unique elements. Its zero value is ready to use.
// Copying a Set by assignment can share its storage. Use Clone for an independent copy.
type Set[T comparable] struct {
	keys map[T]struct{}
}

var itemExists = struct{}{}

// New returns a set containing keys.
func New[T comparable](keys ...T) *Set[T] {
	set := &Set[T]{make(map[T]struct{})}
	set.Add(keys...)
	return set
}

// Add inserts each of keys into the set.
func (st *Set[T]) Add(keys ...T) {
	if st.keys == nil {
		st.keys = make(map[T]struct{})
	}
	for _, k := range keys {
		st.keys[k] = itemExists
	}
}

// Remove deletes each of keys from the set. Absent keys are ignored.
func (st *Set[T]) Remove(keys ...T) {
	for _, k := range keys {
		delete(st.keys, k)
	}
}

// Contains reports whether key is a member of the set.
func (st *Set[T]) Contains(key T) bool {
	_, found := st.keys[key]
	return found
}

// IsEmpty reports whether the set has no elements.
func (st *Set[T]) IsEmpty() bool {
	return len(st.keys) == 0
}

// Len returns the number of elements in the set.
func (st *Set[T]) Len() int {
	return len(st.keys)
}

// Clear removes all elements from the set.
func (st *Set[T]) Clear() {
	clear(st.keys)
}

// String returns a string representation of the set. Element order is unspecified.
func (st *Set[T]) String() string {
	keys := make([]string, 0, len(st.keys))
	for k := range st.keys {
		keys = append(keys, fmt.Sprintf("%v", k))
	}
	return "hashset.Set{" + strings.Join(keys, ", ") + "}"
}

// Clone returns a shallow copy of st with independent storage.
// Elements are copied by assignment.
func (st *Set[T]) Clone() *Set[T] {
	return &Set[T]{maps.Clone(st.keys)}
}
