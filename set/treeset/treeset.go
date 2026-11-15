// Package treeset provides a set ordered by its elements.
package treeset

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/MorePeanuts/gods/maps/treemap"
	"github.com/MorePeanuts/gods/set"
)

var _ set.Interface[int] = (*Set[int])(nil)

// Set is a collection of elements ordered by a comparator. Elements that compare
// equal represent the same member. Use New or NewWithComparator before
// inserting into a Set. A Set must not be copied after first use. Use Clone
// for an independent copy.
type Set[T comparable] struct {
	keys treemap.Map[T, struct{}]
}

// New returns a set containing keys ordered by cmp.Compare.
func New[T cmp.Ordered](keys ...T) *Set[T] {
	return NewWithComparator(cmp.Compare[T], keys...)
}

// NewWithComparator returns a set containing keys ordered by compare.
// The comparator must define a consistent ordering and return zero for
// elements treated as equal. NewWithComparator panics if compare is nil.
func NewWithComparator[T comparable](compare func(a, b T) int, keys ...T) *Set[T] {
	st := &Set[T]{keys: *treemap.NewWithComparator[T, struct{}](compare)}
	st.Add(keys...)
	return st
}

// Add inserts each of keys. An existing element that compares equal is retained.
// Add panics on insertion if no comparator has been initialized.
func (st *Set[T]) Add(keys ...T) {
	for _, key := range keys {
		st.keys.Put(key, struct{}{})
	}
}

// Remove deletes each of keys from the set. Absent keys are ignored.
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

// Clear removes all elements and retains the comparator.
func (st *Set[T]) Clear() {
	st.keys.Clear()
}

// String returns a string representation of the set in comparator order.
func (st *Set[T]) String() string {
	keys := st.keys.Keys()
	values := make([]string, 0, len(keys))
	for _, key := range keys {
		values = append(values, fmt.Sprintf("%v", key))
	}
	return "treeset.Set{" + strings.Join(values, ", ") + "}"
}

// Clone returns a shallow copy of st with independent storage and the same
// comparator. Elements are copied by assignment.
func (st *Set[T]) Clone() *Set[T] {
	return &Set[T]{keys: *st.keys.Clone()}
}
