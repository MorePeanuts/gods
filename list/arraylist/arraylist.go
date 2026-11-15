package arraylist

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/MorePeanuts/gods/list"
)

var _ list.List[int] = (*List[int])(nil)

type List[T comparable] struct {
	elements []T
}

func New[T comparable](vals ...T) *List[T] {
	lst := &List[T]{}
	if len(vals) > 0 {
		lst.Append(vals...)
	}
	return lst
}

func (lst *List[T]) Get(i int) T {
	if !lst.withinRange(i) {
		panic(fmt.Sprintf("arraylist.List.Get: index %d out of range for length %d", i, len(lst.elements)))
	}
	return lst.elements[i]
}

func (lst *List[T]) Set(i int, val T) {
	if !lst.withinRange(i) {
		panic(fmt.Sprintf("arraylist.List.Set: index %d out of range for length %d", i, len(lst.elements)))
	}
	lst.elements[i] = val
}

func (lst *List[T]) Append(vals ...T) {
	lst.elements = append(lst.elements, vals...)
}

func (lst *List[T]) Remove(i int) T {
	if !lst.withinRange(i) {
		panic(fmt.Sprintf("arraylist.List.Remove: index %d out of range for length %d", i, len(lst.elements)))
	}
	element := lst.elements[i]
	lst.elements = slices.Delete(lst.elements, i, i+1)
	return element
}

func (lst *List[T]) Contains(val T) bool {
	return slices.Contains(lst.elements, val)
}

func (lst *List[T]) ContainsAll(vals ...T) bool {
	set := make(map[T]bool)
	for _, val := range vals {
		set[val] = true
	}
	for _, element := range lst.elements {
		if set[element] {
			delete(set, element)
		}
	}
	return len(lst.elements) == 0
}

func (lst *List[T]) SortBy(compare func(a, b T) int) {
	if len(lst.elements) < 2 {
		return
	}
	slices.SortFunc(lst.elements, compare)
}

func (lst *List[T]) SortStableBy(compare func(a, b T) int) {
	if len(lst.elements) < 2 {
		return
	}
	slices.SortStableFunc(lst.elements, compare)
}

func Sort[T cmp.Ordered](lst *List[T]) {
	slices.Sort(lst.elements)
}

func SortStable[T cmp.Ordered](lst *List[T]) {
	slices.SortStableFunc(lst.elements, cmp.Compare[T])
}

func (lst *List[T]) Swap(i, j int) {
	if !lst.withinRange(i) {
		panic(fmt.Sprintf("arraylist.List.Swap: index %d out of range for length %d", i, len(lst.elements)))
	}
	if !lst.withinRange(j) {
		panic(fmt.Sprintf("arraylist.List.Swap: index %d out of range for length %d", j, len(lst.elements)))
	}
	lst.elements[i], lst.elements[j] = lst.elements[j], lst.elements[i]
}

func (lst *List[T]) Insert(i int, vals ...T) {
	if !lst.withinRange(i) && i != len(lst.elements) {
		panic(fmt.Sprintf("arraylist.List.Insert: index %d out of range [0, %d]", i, len(lst.elements)))
	}
	lst.elements = slices.Insert(lst.elements, i, vals...)
}

func (lst *List[T]) IndexOf(val T) int {
	return slices.Index(lst.elements, val)
}

func (lst *List[T]) IsEmpty() bool {
	return len(lst.elements) == 0
}

func (lst *List[T]) Len() int {
	return len(lst.elements)
}

func (lst *List[T]) Clear() {
	clear(lst.elements)
	lst.elements = lst.elements[:0]
}

func (lst *List[T]) String() string {
	vals := make([]string, 0, len(lst.elements))
	for _, val := range lst.elements {
		vals = append(vals, fmt.Sprintf("%v", val))
	}
	return "arraylist.List{" + strings.Join(vals, ", ") + "}"
}

func (lst *List[T]) ShrinkToFit() {
	lst.ShrinkTo(len(lst.elements))
}

func (lst *List[T]) ShrinkTo(n int) {
	target := max(len(lst.elements), n)
	if target < cap(lst.elements) {
		lst.resize(len(lst.elements), target)
	}
}

func (lst *List[T]) resize(length, cap int) {
	newElements := make([]T, length, cap)
	copy(newElements, lst.elements)
	lst.elements = newElements
}

func (lst *List[T]) withinRange(i int) bool {
	return i >= 0 && i < len(lst.elements)
}
