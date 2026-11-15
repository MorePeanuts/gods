// Package linkedhashset provides a hash-based set that preserves insertion order.
package linkedhashset

import (
	"fmt"
	"strings"

	"github.com/MorePeanuts/gods/maps/linkedhashmap"
	"github.com/MorePeanuts/gods/set"
)

var _ set.Interface[int] = (*Set[int])(nil)

// Set is a collection of unique elements in insertion order.
// Its zero value is ready to use. Copying a Set by assignment can share its
// storage. Use Clone for an independent copy.
type Set[T comparable] struct {
	keys linkedhashmap.Map[T, struct{}]
}

// New returns a set containing keys in the order they first appear.
func New[T comparable](keys ...T) *Set[T] {
	st := &Set[T]{}
	st.Add(keys...)
	return st
}

// Add inserts each of keys in argument order. Existing keys keep their positions.
func (st *Set[T]) Add(keys ...T) {
	for _, key := range keys {
		st.keys.Put(key, struct{}{})
	}
}

// Remove deletes each of keys from the set. Absent keys are ignored.
// Adding a removed key again places it at the end.
func (st *Set[T]) Remove(keys ...T) {
	for _, key := range keys {
		st.keys.Remove(key)
	}
}

// Contains reports whether key is a member of the set.
func (st *Set[T]) Contains(key T) bool {
	_, found := st.keys.Get(key)
	return found
}

// IsEmpty reports whether the set has no elements.
func (st *Set[T]) IsEmpty() bool {
	return st.keys.IsEmpty()
}

// Len returns the number of elements in the set.
func (st *Set[T]) Len() int {
	return st.keys.Len()
}

// Clear removes all elements from the set.
func (st *Set[T]) Clear() {
	st.keys.Clear()
}

// String returns a string representation of the set in insertion order.
func (st *Set[T]) String() string {
	keys := st.keys.Keys()
	values := make([]string, 0, len(keys))
	for _, key := range keys {
		values = append(values, fmt.Sprintf("%v", key))
	}
	return "linkedhashset.Set{" + strings.Join(values, ", ") + "}"
}

// Clone returns a shallow copy of st with independent storage and the same
// insertion order. Elements are copied by assignment.
func (st *Set[T]) Clone() *Set[T] {
	return &Set[T]{keys: *st.keys.Clone()}
}
