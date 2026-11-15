// Package container defines operations shared by container types.
package container

// Container is a collection that can be queried, cleared, and formatted.
type Container[T any] interface {
	// IsEmpty reports whether the container has no elements.
	IsEmpty() bool
	// Len returns the number of elements.
	Len() int
	// Clear removes all elements.
	Clear()
	// String returns a string representation of the container.
	String() string
}
