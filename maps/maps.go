// Package maps defines the common interface for maps.
package maps

import "github.com/MorePeanuts/gods/container"

// Interface represents a mutable collection that associates each key with one value.
type Interface[K comparable, V any] interface {
	// Get returns the value associated with key, or the zero value and false
	// if key is absent.
	Get(key K) (V, bool)
	// Put associates val with key. It returns the previous value and true
	// if key was present, or the zero value and false for a new key.
	Put(key K, val V) (old V, replaced bool)
	// Remove deletes key and returns its previous value and true, or the
	// zero value and false if key is absent.
	Remove(key K) (old V, removed bool)
	// Keys returns an independent slice of keys in implementation-defined order.
	Keys() []K

	container.Container[V]
}
