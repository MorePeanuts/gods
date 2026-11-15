// Package set defines the common interface and membership helpers for sets.
package set

import (
	"slices"

	"github.com/MorePeanuts/gods/container"
)

// Interface represents a mutable collection of unique elements.
type Interface[T comparable] interface {
	// Add inserts each of keys into the set.
	Add(keys ...T)
	// Remove deletes each of keys from the set. Absent keys are ignored.
	Remove(keys ...T)
	// Contains reports whether key is a member of the set.
	Contains(key T) bool

	container.Container[T]
}

// ContainsAll reports whether st contains every key. It returns true if keys is empty.
func ContainsAll[T comparable](st Interface[T], keys ...T) bool {
	for _, key := range keys {
		if !st.Contains(key) {
			return false
		}
	}
	return true
}

// ContainsAny reports whether st contains at least one key. It returns false if keys is empty.
func ContainsAny[T comparable](st Interface[T], keys ...T) bool {
	return slices.ContainsFunc(keys, st.Contains)
}
