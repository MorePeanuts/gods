package container

type Container[T any] interface {
	IsEmpty() bool
	Len() int
	Clear()
	String() string
}
