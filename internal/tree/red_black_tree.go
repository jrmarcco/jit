package tree

import (
	"github.com/JrMarcco/jit"
	"github.com/JrMarcco/jit/internal/errs"
)

// color specifies the color of the node
//
// considering memory alignment, there is no difference between in runtime performance between bool and int.
// semantically bool is more in line with the red-black tree's definition.
type color bool

const (
	red   color = false
	black color = true
)

// RBTree is a red-black tree
//  1. the root node is black
//  2. every leaf node is black or nil; that means a leaf node does not store any value
//     2.1 also for space-saving implementation share a black empty node
//  3. any neighboring node (parent and child) cannot be red at the same time
//  4. every path from root to leaf node has the same number of black nodes
type RBTree[K any, V any] struct {
	root *rbNode[K, V]
	size int64
	cmp  jit.Comparator[K]
}

func (rbt *RBTree[K, V]) Size() int64 {
	return rbt.size
}

func (rbt *RBTree[K, V]) Put(key K, val V) error {
	if err := rbt.insertNode(newRbNode(key, val)); err != nil {
		return err
	}

	rbt.size++
	return nil
}

// Del deletes the given key from the tree
func (rbt *RBTree[K, V]) Del(key K) (V, error) {
	node := rbt.findNode(key)
	if node == nil {
		var zero V
		return zero, errs.ErrNodeNotFound
	}

	nodeVal := node.val

	rbt.deleteNode(node)
	rbt.size--
	return nodeVal, nil
}

// Set sets the value of the given key
func (rbt *RBTree[K, V]) Set(key K, val V) error {
	if node := rbt.findNode(key); node != nil {
		node.val = val
		return nil
	}

	return errs.ErrNodeNotFound
}

func (rbt *RBTree[K, V]) Get(key K) (V, error) {
	if node := rbt.findNode(key); node != nil {
		return node.val, nil
	}

	var zero V
	return zero, errs.ErrNodeNotFound
}

// Keys returns the keys of the tree
func (rbt *RBTree[K, V]) Keys() []K {
	keys := make([]K, 0, rbt.size)

	if rbt.root == nil {
		return keys
	}

	rbt.midOrderTraversal(func(node *rbNode[K, V]) {
		keys = append(keys, node.key)
	})
	return keys
}

// Vals returns the values of the tree
func (rbt *RBTree[K, V]) Vals() []V {
	vals := make([]V, 0, rbt.size)

	if rbt.root == nil {
		return vals
	}

	rbt.midOrderTraversal(func(node *rbNode[K, V]) {
		vals = append(vals, node.val)
	})
	return vals
}

// Kvs returns the keys and values of the tree
func (rbt *RBTree[K, V]) Kvs() (keys []K, vals []V) {
	keys = make([]K, 0, rbt.size)
	vals = make([]V, 0, rbt.size)

	if rbt.root == nil {
		return keys, vals
	}

	rbt.midOrderTraversal(func(node *rbNode[K, V]) {
		keys = append(keys, node.key)
		vals = append(vals, node.val)
	})

	return
}

// Iter to traversal all the tree node.
//
// If fn returns false, terminate traversal.
func (rbt *RBTree[K, V]) Iter(visitFunc func(key K, val V) bool) {
	rbt.inOrderTraversal(func(node *rbNode[K, V]) bool {
		return visitFunc(node.key, node.val)
	})
}

func (rbt *RBTree[K, V]) inOrderTraversal(visitFunc func(node *rbNode[K, V]) bool) {
	stack := make([]*rbNode[K, V], 0)
	currNode := rbt.root

	for currNode != nil || len(stack) > 0 {
		for currNode != nil {
			// stack push
			stack = append(stack, currNode)
			currNode = currNode.left
		}

		// stack pop
		currNode = stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if !visitFunc(currNode) {
			break
		}
		currNode = currNode.right
	}
}

// leftRotate left rotate around the node
//
//	     left rotate around the node x
//		 (a / b / r can be a subtree of nil)
//
//		         |                      |
//		         x                      y
//		        / \                    / \
//			   a   y        =>        x   r
//			      / \                / \
//			     b   r              a   b
func (rbt *RBTree[K, V]) leftRotate(x *rbNode[K, V]) {
	if x == nil || x.right == nil {
		// if node x is nil or node x's right is nil, do nothing
		return
	}

	// node y is x's right child
	y := x.right
	// node x's right = node y's left
	x.right = y.left
	// if node y's left is not nil, node y's left's parent = node x
	if y.left != nil {
		y.left.parent = x
	}

	// node y's parent = node x's parent
	y.parent = x.parent

	switch {
	case x.parent == nil:
		// if node x's parent is nil, node x is root, change root to node y
		rbt.root = y
	case x == x.parent.left:
		// if node x is the left child, node y is the left child
		x.parent.left = y
	default:
		// if node x is the right child, node y is the right child
		x.parent.right = y
	}

	// node y's left = node x
	y.left = x
	// node x's parent = node y
	x.parent = y
}

// rightRotate right rotate around the node
//
//	     right rotate around the node x
//		(a / b / r can be a subtree of nil)
//
//		     	|                       |
//		     	x                       y
//		  	   / \                     / \
//		  	  y   r        =>         a   x
//		  	 / \                    	 / \
//		 	a   b                  		b   r
func (rbt *RBTree[K, V]) rightRotate(x *rbNode[K, V]) {
	if x == nil || x.left == nil {
		// if node x is nil or node x's left is nil, do nothing
		return
	}

	// left: node y
	y := x.left
	// node x's left = node y's right
	x.left = y.right
	// if node y's right is not nil, node y's right's parent = node x
	if y.right != nil {
		y.right.parent = x
	}

	// node y's parent = node x's parent
	y.parent = x.parent

	switch {
	case x.parent == nil:
		// if node x's parent is nil, node x is root, change root to node y
		rbt.root = y
	case x == x.parent.right:
		// if node x is the right child, node y is the right child
		x.parent.right = y
	default:
		// if node x is the left child, node y is the left child
		x.parent.left = y
	}

	// node y's right = node x
	y.right = x
	// node x's parent = node y
	x.parent = y
}

// insertNode insert a new node into the tree
// red-black specifies that the inserted node must be red.
func (rbt *RBTree[K, V]) insertNode(node *rbNode[K, V]) error {
	if rbt.root == nil {
		// if the tree is empty, the inserted node is the root
		rbt.root = newRbNode(node.key, node.val)
		rbt.root.setColor(black)
		return nil
	}

	cmp := 0
	parent := &rbNode[K, V]{}

	currNode := rbt.root
	for currNode != nil {
		parent = currNode

		cmp = rbt.cmp(node.key, currNode.key)
		if cmp == 0 {
			return errs.ErrSameRBNode
		}

		if cmp < 0 {
			currNode = currNode.left
		} else {
			currNode = currNode.right
		}
	}

	// the first focus on node is the inserted node
	insertedNode := newRbNode(node.key, node.val)
	insertedNode.parent = parent

	if cmp < 0 {
		parent.left = insertedNode
	} else {
		parent.right = insertedNode
	}

	rbt.fixupInsertion(insertedNode)
	return nil
}

// fixupInsertion ensures the red-black tree properties are maintained after insertion.
// It handles three cases based on the color of the node's uncle:
// 1. Uncle is red
// 2. Uncle is black and its parent is the left child
// 3. Uncle is black and its parent is the right child
func (rbt *RBTree[K, V]) fixupInsertion(node *rbNode[K, V]) {
	for node != nil && node != rbt.root && node.parent.getColor() == red {
		uncle := node.getUncle()
		if uncle.getColor() == red {
			node = rbt.fixupRedUncle(node, uncle)
			continue
		}

		if node.parent == node.getGrandparent().left {
			// case 2: uncle is black, and its parent is left child
			node = rbt.fixupBlackUncleLeftChild(node)
			continue
		}

		// case 3: uncle is black, and its parent is the right child
		node = rbt.fixupBlackUncleRightChild(node)
	}

	// the new inserted node is root,
	// or the new inserted node's parent node is black,
	// no need to fixup
	rbt.root.setColor(black)
}

// fixupRedUncle handles the case where the node's uncle is red.
// It recolors the parent and uncle to black and the grandparent to red,
// then moves the focus to the grandparent.
//
//			    b(c)							           r(c)
//			  /	    \							         /     \
//		  r(b)		r(d)							  b(b)	   b(d)
//		 /	\	    /	\		       =>		     /	\	   /   \
//	  b(z)	r(x)  b(p)  b(q)					  b(z)	r(x)  b(p)  b(q)
//
// 1. change focus on node x's parent node b and uncle node d to black.
func (rbt *RBTree[K, V]) fixupRedUncle(node, uncle *rbNode[K, V]) *rbNode[K, V] {
	grandparent := node.getGrandparent()

	node.parent.setColor(black)
	uncle.setColor(black)
	grandparent.setColor(red)

	return grandparent
}

// fixupBlackUncleLeftChild handles the case where the node's uncle is black,
// and its parent is left child.
//
//			    b(c)							           b(c)
//			  /	    \							         /     \
//		  b(b)		r(d)							  b(b)	   r(d)
//		 /	\	    /  \		       =>		     /	\	   /   \
//	  b(z)	r(y) b(p)  b(q)					      b(z)	r(x) b(p)  b(q)
//			   \							       		/
//		       r(x)								     r(y)
//
// 1. if the focus on node x is the right child of its parent, change the focus on the node to its parent node y.
// 2. left rotate around the focus on node y.
//
//			   b(c)	   	                   b(c)                     r(b)
//			  /	  \	   		              /   \                    /   \
//		  b(b)	   r(d)			      r(b)	   r(d)             b(z)   b(c)
//		 /	\	   /  \	     =>	     /	\	   /  \	     =>            /   \
//	  b(z)	r(x) b(p) b(q)		  b(z)	b(x) b(p) b(q)              b(x)   r(d)
//			/                           /                            /     /  \
//		  r(y)                        r(y)                        r(y)    b(p)  b(q)
//
// 3. change focus on node y's parent to black.
// 4. change focus on node y's grandparent to red.
// 5. right rotate around the focus on node y's grandparent b.
func (rbt *RBTree[K, V]) fixupBlackUncleLeftChild(node *rbNode[K, V]) *rbNode[K, V] {
	if node == node.parent.right {
		node = node.parent
		rbt.leftRotate(node)
	}

	node.parent.setColor(black)
	node.getGrandparent().setColor(red)
	rbt.rightRotate(node.getGrandparent())

	return node.parent
}

// fixupBlackUncleRightChild handles the case where the node's uncle is black,
// and its parent is right child.
//
//			    b(c)							           b(c)
//			  /	    \							         /     \
//		  r(b)		b(d)							  r(b)	   b(d)
//		 /	\	    /  \		       =>		     /	\	   /  \
//	  b(p)	b(q) b(z) r(y)					      b(p)	b(q) b(z) r(x)
//			          /                                             \
//		             r(x)                                           r(y)
//
// 1. if the focus on node x is the left child of its parent, change the focus on the node to its parent node y.
// 2. right rotate around the focus on node y.
//
//			   b(c)	   	                   b(c)                         r(d)
//			  /	  \	   		              /   \                        /   \
//		  r(b)	   b(d)			      r(b)	   r(d)                 b(c)   b(x)
//		 /	\	   /  \	     =>	     /	\	   /  \	       =>       /  \     \
//	  b(p)	b(q) b(z) r(x)		  b(p)	b(q) b(z) b(x)            r(b) b(z)  r(y)
//			          	\                           \             /  \
//		             	r(y)						r(y)	   b(p)  b(q)
//
// 3. change focus on node y's parent to black.
// 4. change focus on node y's grandparent to red.
// 5. left rotate around the focus on node y's grandparent b.
func (rbt *RBTree[K, V]) fixupBlackUncleRightChild(node *rbNode[K, V]) *rbNode[K, V] {
	if node == node.parent.left {
		node = node.parent
		rbt.rightRotate(node)
	}

	node.parent.setColor(black)
	node.getGrandparent().setColor(red)
	rbt.leftRotate(node.getGrandparent())

	return node
}

// findSuccessor find the successor of the given node.
//
//		             N
//			        / \
//		           L  R
//		             / \
//		           RL  RR
//		          /
//		        RLL
//	1. If the node has a right child,
//	   	the successor is the leftmost node of the right subtree.
//		From the right subtree of the deleted node(N),
//		traverse all the way to the leftmost node(RLL).(the node with the smallest key in the right subtree)
//
//	            P
//	             \
//		          N
//		         /
//		        L
//	2. If the node has no right child,
//		find the deleted node(N)'s parent node.
//		if the deleted node is the left child of its parent, the successor is the parent node.
//		Otherwise, backtrack up the parent node until find the first ancestor node that is the left child of its parent.
//		The ancestor node's parent is the successor.
func (rbt *RBTree[K, V]) findSuccessor(node *rbNode[K, V]) *rbNode[K, V] {
	if node == nil {
		return nil
	}

	if node.right != nil {
		// if the node has a right child, the successor is the leftmost node of the right subtree
		curr := node.right
		for curr.left != nil {
			curr = curr.left
		}
		return curr
	}

	// the node has no right child
	parent := node.parent
	curr := node

	for parent != nil && curr == parent.right {
		// keep moving up until the node is the left child of its parent
		curr = parent
		parent = parent.parent
	}

	return parent
}

// deleteNode delete the given node from the tree.
func (rbt *RBTree[K, V]) deleteNode(node *rbNode[K, V]) {
	deletedNode := node

	if deletedNode.left != nil && deletedNode.right != nil {
		// if the deleted node has both left and right child,
		// find the successor of the deleted node.
		successor := rbt.findSuccessor(deletedNode)

		// copy the successor's key and val to the deleted node
		deletedNode.key = successor.key
		deletedNode.val = successor.val

		// replace the deleted node with the successor,
		// and the successor will be deleted in the next step
		deletedNode = successor
	}

	var replacement *rbNode[K, V]

	// now deletedNode is the successor of the original deleted node
	// if the deleted node has left child, the replacement is the left child
	// otherwise, the replacement is the right child
	if deletedNode.left != nil {
		replacement = deletedNode.left
	} else {
		replacement = deletedNode.right
	}

	if replacement != nil {
		// replace the deleted node with the replacement
		replacement.parent = deletedNode.parent
		switch {
		case deletedNode.parent == nil:
			// if the deleted node is the root, replace the root with the replacement
			rbt.root = replacement
		case deletedNode == deletedNode.parent.left:
			// if the deleted node is the left child of its parent, replace the deleted node with the replacement
			deletedNode.parent.left = replacement
		default:
			// if the deleted node is the right child of its parent, replace the deleted node with the replacement
			deletedNode.parent.right = replacement
		}

		deletedNode.left = nil
		deletedNode.right = nil
		deletedNode.parent = nil

		if deletedNode.getColor() == black {
			rbt.deletionFixup(replacement)
		}

		return
	}

	if deletedNode.parent == nil {
		rbt.root = nil
		return
	}

	if deletedNode.getColor() == black {
		rbt.deletionFixup(deletedNode)
	}

	if deletedNode.parent != nil {
		switch deletedNode {
		case deletedNode.parent.left:
			deletedNode.parent.left = nil
		case deletedNode.parent.right:
			deletedNode.parent.right = nil
		}
		deletedNode.parent = nil
	}
}

// deletionFixup fixes the red-black tree properties after deletion.
func (rbt *RBTree[K, V]) deletionFixup(node *rbNode[K, V]) {
	for node != rbt.root && node.getColor() == black {
		if node == node.parent.left {
			node = rbt.deletionFixupLeftChild(node)
			continue
		}
		node = rbt.deletionFixupRightChild(node)
	}

	node.setColor(black)
}

// deletionFixupLeftChild fixes the red-black tree properties after deletion for the left child.
func (rbt *RBTree[K, V]) deletionFixupLeftChild(node *rbNode[K, V]) *rbNode[K, V] {
	// right child of the parent node
	rcNode := node.parent.right

	if rcNode.getColor() == red {
		rcNode.setColor(black)
		rcNode.parent.setColor(red)
		rbt.leftRotate(node.parent)
		rcNode = node.parent.right
	}

	if rcNode.left.getColor() == black && rcNode.right.getColor() == black {
		rcNode.setColor(red)
		return node.parent
	}

	if rcNode.right.getColor() == black {
		rcNode.left.setColor(black)
		rcNode.setColor(red)
		rbt.rightRotate(rcNode)
		rcNode = node.parent.right
	}

	rcNode.setColor(node.parent.getColor())
	node.parent.setColor(black)
	rcNode.right.setColor(black)
	rbt.leftRotate(node.parent)

	return rbt.root
}

// deletionFixupRightChild fixes the red-black tree properties after deletion for the right child.
func (rbt *RBTree[K, V]) deletionFixupRightChild(node *rbNode[K, V]) *rbNode[K, V] {
	// left child of the parent node
	lcNode := node.parent.left

	if lcNode.getColor() == red {
		lcNode.setColor(black)
		node.parent.setColor(red)
		rbt.rightRotate(node.parent)
		lcNode = node.getSibling()
	}

	if lcNode.right.getColor() == black && lcNode.left.getColor() == black {
		lcNode.setColor(red)
		return node.parent
	}

	if lcNode.left.getColor() == black {
		lcNode.right.setColor(black)
		lcNode.setColor(red)
		rbt.leftRotate(lcNode)
		lcNode = node.parent.left
	}

	lcNode.setColor(node.parent.getColor())
	node.parent.setColor(black)
	lcNode.left.setColor(black)
	rbt.rightRotate(node.parent)

	return rbt.root
}

func (rbt *RBTree[K, V]) findNode(key K) *rbNode[K, V] {
	node := rbt.root
	for node != nil {
		cmp := rbt.cmp(key, node.key)
		if cmp == 0 {
			return node
		}

		if cmp < 0 {
			node = node.left
		} else {
			node = node.right
		}
	}

	return nil
}

func (rbt *RBTree[K, V]) midOrderTraversal(visitFn func(node *rbNode[K, V])) {
	stack := make([]*rbNode[K, V], 0, rbt.size)

	curr := rbt.root
	for curr != nil || len(stack) > 0 {
		for curr != nil {
			stack = append(stack, curr)
			curr = curr.left
		}

		curr = stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		visitFn(curr)

		curr = curr.right
	}
}

func NewRBTree[K any, V any](cmp jit.Comparator[K]) *RBTree[K, V] {
	return &RBTree[K, V]{
		root: nil,
		size: 0,
		cmp:  cmp,
	}
}

// rbNode is a node of a red-black tree
// consider memory alignment, put color at the end of the struct
type rbNode[K any, V any] struct {
	parent *rbNode[K, V]
	left   *rbNode[K, V]
	right  *rbNode[K, V]
	key    K
	val    V
	color  color
}

func (rbn *rbNode[K, V]) getColor() color {
	if rbn == nil {
		return black
	}
	return rbn.color
}

func (rbn *rbNode[K, V]) setColor(color color) {
	if rbn == nil {
		return
	}
	rbn.color = color
}

// getGrandparent get the grandparent node
func (rbn *rbNode[K, V]) getGrandparent() *rbNode[K, V] {
	if rbn == nil || rbn.parent == nil {
		return nil
	}
	return rbn.parent.parent
}

// getSibling get the sibling(brother) node
func (rbn *rbNode[K, V]) getSibling() *rbNode[K, V] {
	if rbn == nil || rbn.parent == nil {
		return nil
	}
	if rbn == rbn.parent.left {
		return rbn.parent.right
	}
	return rbn.parent.left
}

// findUncle find the uncle node of the given node
func (rbn *rbNode[K, V]) getUncle() *rbNode[K, V] {
	if rbn == nil {
		return nil
	}

	return rbn.parent.getSibling()
}

// newRbNode create a new node
// the new node is red before insert fixup
func newRbNode[K any, V any](key K, val V) *rbNode[K, V] {
	return &rbNode[K, V]{
		key:    key,
		val:    val,
		color:  red,
		parent: nil,
		left:   nil,
		right:  nil,
	}
}
