// Package rbtree provides a left-leaning red-black search tree.
package rbtree

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/MorePeanuts/gods/tree"
)

var _ tree.SearchTree[int, string] = (*Tree[int, string])(nil)

// Tree stores key-value pairs ordered by a comparator.
// Keys that compare equal represent the same entry. An empty zero value can
// be queried, but insertion requires New or NewWithComparator. A Tree must
// not be copied after first use. Use Clone for an independent copy.
type Tree[K any, V any] struct {
	root    *node[K, V]
	length  int
	compare func(a, b K) int
}

// New returns an empty tree ordered by cmp.Compare.
func New[K cmp.Ordered, V any]() *Tree[K, V] {
	return NewWithComparator[K, V](cmp.Compare[K])
}

// NewWithComparator returns an empty tree ordered by compare. The comparator
// must define a consistent ordering and return zero for keys treated as equal.
// NewWithComparator panics if compare is nil.
func NewWithComparator[K any, V any](compare func(a, b K) int) *Tree[K, V] {
	if compare == nil {
		panic("rbtree: nil comparator")
	}
	return &Tree[K, V]{compare: compare}
}

func (tr *Tree[K, V]) find(key K) *node[K, V] {
	for n := tr.root; n != nil; {
		switch order := tr.compare(key, n.key); {
		case order < 0:
			n = n.left
		case order > 0:
			n = n.right
		default:
			return n
		}
	}
	return nil
}

// Get returns the value associated with key and whether key is present.
// If key is absent, it returns the zero value and false.
func (tr *Tree[K, V]) Get(key K) (V, bool) {
	if n := tr.find(key); n != nil {
		return n.value, true
	}
	return *new(V), false
}

// Put associates val with key. It returns the previous value and true if key
// was present, or the zero value and false for a new key. An existing entry
// retains its original key. Put panics if no comparator has been initialized.
func (tr *Tree[K, V]) Put(key K, val V) (old V, replaced bool) {
	if tr.compare == nil {
		panic("rbtree: initialize with New or NewWithComparator before inserting")
	}
	tr.root, old, replaced = tr.insert(tr.root, key, val)
	tr.root.red = false
	if !replaced {
		tr.length++
	}
	return old, replaced
}

// Remove deletes key and returns its previous value and true.
// If key is absent, it returns the zero value and false.
func (tr *Tree[K, V]) Remove(key K) (old V, removed bool) {
	old, removed = tr.Get(key)
	if !removed {
		return old, false
	}

	// Start deletion with a red root when both children are black so
	// the descent can move a red link toward the node being removed.
	if !isRed(tr.root.left) && !isRed(tr.root.right) {
		tr.root.red = true
	}
	tr.root = tr.remove(tr.root, key)
	if tr.root != nil {
		tr.root.red = false
	}
	tr.length--
	return old, true
}

// Ascend calls yield for each entry in comparator order until yield returns
// false. The callback must not modify the tree.
func (tr *Tree[K, V]) Ascend(yield func(K, V) bool) {
	tr.root.ascend(yield)
}

// First returns the first entry in comparator order, or zero values and false
// if the tree is empty.
func (tr *Tree[K, V]) First() (K, V, bool) {
	n := tr.root
	for n != nil && n.left != nil {
		n = n.left
	}
	return n.keyValue()
}

// Last returns the last entry in comparator order, or zero values and false
// if the tree is empty.
func (tr *Tree[K, V]) Last() (K, V, bool) {
	n := tr.root
	for n != nil && n.right != nil {
		n = n.right
	}
	return n.keyValue()
}

// Floor returns the last entry whose key is at or before key in comparator
// order, or zero values and false if no such entry exists.
func (tr *Tree[K, V]) Floor(key K) (K, V, bool) {
	return tr.predecessor(key, true).keyValue()
}

// Ceiling returns the first entry whose key is at or after key in comparator
// order, or zero values and false if no such entry exists.
func (tr *Tree[K, V]) Ceiling(key K) (K, V, bool) {
	return tr.successor(key, true).keyValue()
}

// Lower returns the last entry whose key is strictly before key in comparator
// order, or zero values and false if no such entry exists.
func (tr *Tree[K, V]) Lower(key K) (K, V, bool) {
	return tr.predecessor(key, false).keyValue()
}

// Higher returns the first entry whose key is strictly after key in comparator
// order, or zero values and false if no such entry exists.
func (tr *Tree[K, V]) Higher(key K) (K, V, bool) {
	return tr.successor(key, false).keyValue()
}

func (tr *Tree[K, V]) predecessor(key K, inclusive bool) *node[K, V] {
	var best *node[K, V]
	for n := tr.root; n != nil; {
		order := tr.compare(key, n.key)
		if order == 0 && inclusive {
			return n
		}
		if order > 0 {
			best = n
			n = n.right
		} else {
			n = n.left
		}
	}
	return best
}

func (tr *Tree[K, V]) successor(key K, inclusive bool) *node[K, V] {
	var best *node[K, V]
	for n := tr.root; n != nil; {
		order := tr.compare(key, n.key)
		if order == 0 && inclusive {
			return n
		}
		if order < 0 {
			best = n
			n = n.left
		} else {
			n = n.right
		}
	}
	return best
}

func (n *node[K, V]) keyValue() (K, V, bool) {
	if n == nil {
		return *new(K), *new(V), false
	}
	return n.key, n.value, true
}

// AscendRange calls yield for entries in the half-open interval [from, to)
// in comparator order until yield returns false. It yields nothing if from
// compares equal to or after to. The callback must not modify the tree.
// The traversal takes O(log n + k) time for k visited entries.
func (tr *Tree[K, V]) AscendRange(from, to K, yield func(K, V) bool) {
	if tr.root == nil || tr.compare(from, to) >= 0 {
		return
	}

	var visit func(*node[K, V]) bool
	visit = func(n *node[K, V]) bool {
		if n == nil {
			return true
		}
		lower, upper := tr.compare(n.key, from), tr.compare(n.key, to)
		if lower > 0 && !visit(n.left) {
			return false
		}
		if lower >= 0 && upper < 0 && !yield(n.key, n.value) {
			return false
		}
		return upper >= 0 || visit(n.right)
	}
	visit(tr.root)
}

// IsEmpty reports whether the tree has no entries.
func (tr *Tree[K, V]) IsEmpty() bool {
	return tr.length == 0
}

// Len returns the number of entries in the tree.
func (tr *Tree[K, V]) Len() int {
	return tr.length
}

// Clear removes all entries and retains the comparator.
func (tr *Tree[K, V]) Clear() {
	tr.root = nil
	tr.length = 0
}

// String renders the structure of the tree, with L and R marking left and
// right children and [B] and [R] marking black and red nodes.
func (tr *Tree[K, V]) String() string {
	if tr.root == nil {
		return "rbtree.Tree{}"
	}

	var out strings.Builder
	out.WriteString("rbtree.Tree{\n")
	var writeNode func(*node[K, V], string, string, string)
	writeNode = func(n *node[K, V], prefix, branch, side string) {
		color := "B"
		if n.red {
			color = "R"
		}
		fmt.Fprintf(&out, "%s%s%s[%s] %v: %v\n", prefix, branch, side, color, n.key, n.value)

		switch branch {
		case "├── ":
			prefix += "│   "
		case "└── ":
			prefix += "    "
		}
		if n.left != nil {
			leftBranch := "└── "
			if n.right != nil {
				leftBranch = "├── "
			}
			writeNode(n.left, prefix, leftBranch, "L ")
		}
		if n.right != nil {
			writeNode(n.right, prefix, "└── ", "R ")
		}
	}
	writeNode(tr.root, "  ", "", "")
	out.WriteByte('}')
	return out.String()
}

// Clone returns a shallow copy with independent nodes and the same comparator.
// Keys and values are copied by assignment.
func (tr *Tree[K, V]) Clone() *Tree[K, V] {
	return &Tree[K, V]{root: tr.root.clone(), length: tr.length, compare: tr.compare}
}
