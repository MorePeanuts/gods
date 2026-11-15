package linkedlist

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/MorePeanuts/gods/list"
)

var _ list.List[int] = (*List[int])(nil)

type List[T comparable] struct {
	head   *node[T]
	length int
}

type node[T comparable] struct {
	val  T
	prev *node[T]
	next *node[T]
}

func New[T comparable](vals ...T) *List[T] {
	head := &node[T]{}
	head.next = head
	head.prev = head
	lst := &List[T]{head, 0}
	lst.Append(vals...)
	return lst
}

func (lst *List[T]) Get(i int) T {
	if !lst.withinRange(i) {
		panic(fmt.Sprintf("linkedlist.List.Get: index %d out of range for length %d", i, lst.length))
	}
	p := lst.head
	if lst.length-i > i+1 {
		for range i + 1 {
			p = p.next
		}
	} else {
		for range lst.length - i {
			p = p.prev
		}
	}
	return p.val
}

func (lst *List[T]) Set(i int, val T) {
	if !lst.withinRange(i) {
		panic(fmt.Sprintf("linkedlist.List.Set: index %d out of range for length %d", i, lst.length))
	}
	p := lst.head
	if lst.length-i > i+1 {
		for range i + 1 {
			p = p.next
		}
	} else {
		for range lst.length - i {
			p = p.prev
		}
	}
	p.val = val
}

func (lst *List[T]) Append(vals ...T) {
	for _, val := range vals {
		newNode := &node[T]{val, lst.head.prev, lst.head}
		lst.head.prev.next = newNode
		lst.head.prev = newNode
		lst.length++
	}
}

func (lst *List[T]) Remove(i int) T {
	if !lst.withinRange(i) {
		panic(fmt.Sprintf("linkedlist.List.Remove: index %d out of range for length %d", i, lst.length))
	}
	p := lst.head
	if lst.length-i > i+1 {
		for range i + 1 {
			p = p.next
		}
	} else {
		for range lst.length - i {
			p = p.prev
		}
	}
	p.next.prev = p.prev
	p.prev.next = p.next
	lst.length--
	return p.val
}

func (lst *List[T]) Contains(val T) bool {
	for p := lst.head.next; p != lst.head; p = p.next {
		if p.val == val {
			return true
		}
	}
	return false
}

func (lst *List[T]) ContainsAll(vals ...T) bool {
	set := make(map[T]bool)
	for _, val := range vals {
		set[val] = true
	}
	for p := lst.head.next; p != lst.head; p = p.next {
		if set[p.val] {
			delete(set, p.val)
		}
	}
	return len(set) == 0
}

func (lst *List[T]) SortBy(compare func(a, b T) int) {
	vals := make([]T, 0, lst.length)
	for p := lst.head.next; p != lst.head; p = p.next {
		vals = append(vals, p.val)
	}
	slices.SortFunc(vals, compare)
	lst.Clear()
	lst.Append(vals...)
}

func (lst *List[T]) SortStableBy(compare func(a, b T) int) {
	vals := make([]T, 0, lst.length)
	for p := lst.head.next; p != lst.head; p = p.next {
		vals = append(vals, p.val)
	}
	slices.SortStableFunc(vals, compare)
	lst.Clear()
	lst.Append(vals...)
}

func Sort[T cmp.Ordered](lst *List[T]) {
	lst.SortBy(cmp.Compare[T])
}

func SortStable[T cmp.Ordered](lst *List[T]) {
	lst.SortStableBy(cmp.Compare[T])
}

func (lst *List[T]) Swap(i, j int) {
	if !lst.withinRange(i) {
		panic(fmt.Sprintf("linkedlist.List.Swap: index %d out of range for length %d", i, lst.length))
	}
	if !lst.withinRange(j) {
		panic(fmt.Sprintf("linkedlist.List.Swap: index %d out of range for length %d", j, lst.length))
	}
	if i == j {
		return
	}
	var pi, pj *node[T]
	for k, p := 0, lst.head.next; pi == nil || pj == nil; k, p = k+1, p.next {
		switch k {
		case i:
			pi = p
		case j:
			pj = p
		}
	}
	pi.val, pj.val = pj.val, pi.val
}

func (lst *List[T]) Insert(i int, vals ...T) {
	if !lst.withinRange(i) && i != lst.length {
		panic(fmt.Sprintf("linkedlist.List.Insert: index %d out of range [0, %d]", i, lst.length))
	}
	p := lst.head
	if lst.length-i+1 > i {
		for range i {
			p = p.next
		}
	} else {
		for range lst.length - i + 1 {
			p = p.prev
		}
	}
	for _, val := range vals {
		newNode := &node[T]{val, p, p.next}
		newNode.prev.next = newNode
		newNode.next.prev = newNode
		p = newNode
		lst.length++
	}
}

func (lst *List[T]) IndexOf(val T) int {
	for i, p := 0, lst.head.next; p != lst.head; i, p = i+1, p.next {
		if val == p.val {
			return i
		}
	}
	return -1
}

func (lst *List[T]) IsEmpty() bool {
	return lst.length == 0
}

func (lst *List[T]) Len() int {
	return lst.length
}

func (lst *List[T]) Clear() {
	lst.head.next = lst.head
	lst.head.prev = lst.head
	lst.length = 0
}

func (lst *List[T]) String() string {
	vals := make([]string, 0)
	for n := lst.head.next; n != lst.head; n = n.next {
		vals = append(vals, fmt.Sprintf("%v", n.val))
	}
	return "linkedlist.List{" + strings.Join(vals, ", ") + "}"
}

func (lst *List[T]) Prepend(vals ...T) {
	lst.Insert(0, vals...)
}

func (lst *List[T]) withinRange(i int) bool {
	return i >= 0 && i < lst.length
}
