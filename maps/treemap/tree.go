package treemap

// red is the color of the link from the parent. Missing links are black.
type node[K comparable, V any] struct {
	key         K
	value       V
	left, right *node[K, V]
	red         bool
}

func isRed[K comparable, V any](n *node[K, V]) bool {
	return n != nil && n.red
}

func (n *node[K, V]) rotateLeft() *node[K, V] {
	child := n.right
	n.right = child.left
	child.left = n
	child.red, n.red = n.red, true
	return child
}

func (n *node[K, V]) rotateRight() *node[K, V] {
	child := n.left
	n.left = child.right
	child.right = n
	child.red, n.red = n.red, true
	return child
}

func (n *node[K, V]) flipColors() {
	n.red = !n.red
	n.left.red = !n.left.red
	n.right.red = !n.right.red
}

// balance restores left-leaning links after a recursive change.
func (n *node[K, V]) balance() *node[K, V] {
	if isRed(n.right) && !isRed(n.left) {
		n = n.rotateLeft()
	}
	if isRed(n.left) && isRed(n.left.left) {
		n = n.rotateRight()
	}
	if isRed(n.left) && isRed(n.right) {
		n.flipColors()
	}
	return n
}

func (mp *Map[K, V]) insert(n *node[K, V], key K, val V) (root *node[K, V], old V, replaced bool) {
	if n == nil {
		return &node[K, V]{key: key, value: val, red: true}, old, false
	}

	switch order := mp.compare(key, n.key); {
	case order < 0:
		n.left, old, replaced = mp.insert(n.left, key, val)
	case order > 0:
		n.right, old, replaced = mp.insert(n.right, key, val)
	default:
		old = n.value
		n.value = val
		return n, old, true
	}
	return n.balance(), old, replaced
}

// moveRedLeft prepares a descent into a black left child with a black left link.
func (n *node[K, V]) moveRedLeft() *node[K, V] {
	n.flipColors()
	if isRed(n.right.left) {
		n.right = n.right.rotateRight()
		n = n.rotateLeft()
		n.flipColors()
	}
	return n
}

// moveRedRight prepares a descent into a black right child with a black left link.
func (n *node[K, V]) moveRedRight() *node[K, V] {
	n.flipColors()
	if isRed(n.left.left) {
		n = n.rotateRight()
		n.flipColors()
	}
	return n
}

// remove descends only into subtrees known to contain key. It moves a red link
// down as needed so the final removal does not reduce the black height.
func (mp *Map[K, V]) remove(n *node[K, V], key K) *node[K, V] {
	if mp.compare(key, n.key) < 0 {
		if !isRed(n.left) && !isRed(n.left.left) {
			n = n.moveRedLeft()
		}
		n.left = mp.remove(n.left, key)
		return n.balance()
	}

	if isRed(n.left) {
		n = n.rotateRight()
	}
	if mp.compare(key, n.key) == 0 && n.right == nil {
		return nil
	}
	if !isRed(n.right) && !isRed(n.right.left) {
		n = n.moveRedRight()
	}

	if mp.compare(key, n.key) == 0 {
		successor := n.right
		for successor.left != nil {
			successor = successor.left
		}
		n.key, n.value = successor.key, successor.value
		n.right = n.right.removeMin()
	} else {
		n.right = mp.remove(n.right, key)
	}
	return n.balance()
}

func (n *node[K, V]) removeMin() *node[K, V] {
	if n.left == nil {
		return nil
	}
	if !isRed(n.left) && !isRed(n.left.left) {
		n = n.moveRedLeft()
	}
	n.left = n.left.removeMin()
	return n.balance()
}

func (n *node[K, V]) visit(f func(K, V)) {
	if n == nil {
		return
	}
	n.left.visit(f)
	f(n.key, n.value)
	n.right.visit(f)
}

func (n *node[K, V]) clone() *node[K, V] {
	if n == nil {
		return nil
	}
	return &node[K, V]{
		key: n.key, value: n.value, red: n.red,
		left: n.left.clone(), right: n.right.clone(),
	}
}
