package tree

import (
	"math/bits"

	"github.com/jrmarcco/jit"
	"github.com/jrmarcco/jit/internal/errs"
)

// color 表示节点颜色。
//
// 从内存对齐角度看使用 bool 或 int 在运行时差异很小。
// 这里选择 bool，其语义上更贴近红黑树定义。
type color bool

const (
	red   color = false
	black color = true
)

// RBTree 是红黑树实现：
//  1. 根节点为黑色。
//  2. 所有叶子 ( 空子树 ) 视为黑色。
//     本实现使用共享的黑色哨兵节点。
//  3. 任意相邻节点 ( 父子 ) 不能同时为红色。
//  4. 从任意节点到其所有叶子的路径，黑色节点数量相同。
type RBTree[K any, V any] struct {
	root    *rbNode[K, V]
	nilNode *rbNode[K, V]
	size    int64
	cmp     jit.Comparator[K]
}

func NewRBTree[K any, V any](cmp jit.Comparator[K]) *RBTree[K, V] {
	nilNode := &rbNode[K, V]{
		color: black,
		isNil: true,
	}

	return &RBTree[K, V]{
		root:    nil,
		nilNode: nilNode,
		size:    0,
		cmp:     cmp,
	}
}

func (rbt *RBTree[K, V]) Size() int64 {
	return rbt.size
}

func (rbt *RBTree[K, V]) Put(key K, val V) error {
	inserted, _, err := rbt.upsertNode(key, val, true, false)
	if err != nil {
		return err
	}

	if inserted {
		rbt.size++
	}
	return nil
}

// Upsert 在 key 不存在时插入，存在时覆盖更新。
func (rbt *RBTree[K, V]) Upsert(key K, val V) {
	inserted, _, _ := rbt.upsertNode(key, val, true, true)
	if inserted {
		rbt.size++
	}
}

// Del 删除给定 key 并返回被删除的值。
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

// Set 仅在 key 存在时更新值。
// 不存在时返回 ErrNodeNotFound。
func (rbt *RBTree[K, V]) Set(key K, val V) error {
	_, updated, err := rbt.upsertNode(key, val, false, true)
	if err != nil {
		return err
	}
	if !updated {
		return errs.ErrNodeNotFound
	}
	return nil
}

// Get 返回给定 key 对应的值。
// 不存在时返回 ErrNodeNotFound。
func (rbt *RBTree[K, V]) Get(key K) (V, error) {
	if node := rbt.findNode(key); node != nil {
		return node.val, nil
	}

	var zero V
	return zero, errs.ErrNodeNotFound
}

// Keys 返回按中序排列的全部 key。
func (rbt *RBTree[K, V]) Keys() []K {
	keys := make([]K, 0, rbt.size)

	if rbt.root == nil {
		return keys
	}

	rbt.inOrderTraversalAll(func(node *rbNode[K, V]) {
		keys = append(keys, node.key)
	})
	return keys
}

// Vals 返回按 key 有序对应的全部 value。
func (rbt *RBTree[K, V]) Vals() []V {
	vals := make([]V, 0, rbt.size)

	if rbt.root == nil {
		return vals
	}

	rbt.inOrderTraversalAll(func(node *rbNode[K, V]) {
		vals = append(vals, node.val)
	})
	return vals
}

// Kvs 返回按中序顺序对齐的 keys 与 values。
func (rbt *RBTree[K, V]) Kvs() (keys []K, vals []V) {
	keys = make([]K, 0, rbt.size)
	vals = make([]V, 0, rbt.size)

	if rbt.root == nil {
		return keys, vals
	}

	rbt.inOrderTraversalAll(func(node *rbNode[K, V]) {
		keys = append(keys, node.key)
		vals = append(vals, node.val)
	})

	return
}

// Iter 中序遍历整棵树 ( 可中断 )。
//
// 当 visitFunc 返回 false 时提前终止遍历。
func (rbt *RBTree[K, V]) Iter(visitFunc func(key K, val V) bool) {
	rbt.inOrderTraversalBreakable(func(node *rbNode[K, V]) bool {
		return visitFunc(node.key, node.val)
	})
}

// leftRotate 以 x 为支点做左旋。
//
//	         ( a / b / r 可以是空子树 )
//	         |                      |
//	         x                      y
//	        / \                    / \
//		   a   y        =>        x   r
//		      / \                / \
//		     b   r              a   b
func (rbt *RBTree[K, V]) leftRotate(x *rbNode[K, V]) {
	if x == nil || x.right == nil || x.right == rbt.nilNode {
		// x 不可左旋：右孩子不存在（或为哨兵）
		return
	}

	// y = x.right 将 y 提升为新的子树根。
	y := x.right

	// x.right 挂到 y.left。
	x.right = y.left

	// 修复被移动子树的父指针。
	if y.left != nil && y.left != rbt.nilNode {
		y.left.parent = x
	}

	// y 接管 x 原来的父子关系。
	y.parent = x.parent

	switch {
	case x.parent == nil:
		// x 原本是根旋转后 y 成为新根。
		rbt.root = y
	case x == x.parent.left:
		// x 是左孩子。
		x.parent.left = y
	default:
		// x 是右孩子。
		x.parent.right = y
	}

	// x 下沉为 y.left。
	y.left = x
	x.parent = y
}

// rightRotate 以 x 为支点做右旋。
//
//	        ( a / b / r 可以是空子树 )
//	    	|                      |
//	    	x                      y
//	 	   / \                    / \
//	 	  y   r        =>        a   x
//	 	 / \                    	/ \
//		a   b                  	   b   r
func (rbt *RBTree[K, V]) rightRotate(x *rbNode[K, V]) {
	if x == nil || x.left == nil || x.left == rbt.nilNode {
		// x 不可右旋 ( 左孩子不存在 或 为哨兵 )。
		return
	}

	// y = x.left 将 y 提升为新的子树根。
	y := x.left

	// x.left 挂到 y.right。
	x.left = y.right

	// 修复被移动子树的父指针。
	if y.right != nil && y.right != rbt.nilNode {
		y.right.parent = x
	}

	// y 接管 x 原来的父子关系。
	y.parent = x.parent

	switch {
	case x.parent == nil:
		// x 原本是根旋转后 y 成为新根。
		rbt.root = y
	case x == x.parent.right:
		// x 是右孩子。
		x.parent.right = y
	default:
		// x 是左孩子。
		x.parent.left = y
	}

	// x 下沉为 y.right。
	y.right = x
	x.parent = y
}

// upsertNode 根据标志位执行“插入/更新”：
//   - allowInsert=true  允许不存在时插入。
//   - allowUpdate=true  允许存在时更新。
//
// 返回值含义：
//   - inserted: 本次是否发生了新插入。
//   - updated:  本次是否发生了更新。
func (rbt *RBTree[K, V]) upsertNode(
	key K, val V,
	allowInsert, allowUpdate bool,
) (inserted, updated bool, err error) {
	cmpFunc := rbt.cmp

	if rbt.root == nil {
		if allowInsert {
			insertedNode := rbt.newNode(key, val)
			// 空树插入。
			// 新节点直接作为根并染黑。
			rbt.root = insertedNode
			rbt.root.setColor(black)
			return true, false, nil
		}
		return false, false, nil
	}

	cmp := 0
	var parent *rbNode[K, V]

	currNode := rbt.root
	for currNode != nil && currNode != rbt.nilNode {
		parent = currNode

		cmp = cmpFunc(key, currNode.key)
		if cmp == 0 {
			// 命中已有 key。
			// 根据 allowUpdate 决定覆盖或报重复错误。
			if allowUpdate {
				currNode.val = val
				return false, true, nil
			}
			return false, false, errs.ErrSameRBNode
		}

		if cmp < 0 {
			currNode = currNode.left
		} else {
			currNode = currNode.right
		}
	}

	if !allowInsert {
		// 未命中且不允许插入。
		// 调用方按语义自行处理 ( 如 Set 返回 not found )。
		return false, false, nil
	}

	insertedNode := rbt.newNode(key, val)
	insertedNode.parent = parent

	if cmp < 0 {
		parent.left = insertedNode
	} else {
		parent.right = insertedNode
	}

	rbt.fixupInsertion(insertedNode)
	return true, false, nil
}

// fixupInsertion 在插入后修复红黑树性质。
//
// 核心分 3 类：
//  1. 叔叔节点为红   ：父叔染黑，祖父染红，继续向上修复。
//  2. 叔叔黑 + 父在左：按 LR/LL 做旋转与重染色。
//  3. 叔叔黑 + 父在右：按 RL/RR 做旋转与重染色。
func (rbt *RBTree[K, V]) fixupInsertion(node *rbNode[K, V]) {
	for node != nil && node != rbt.root && node.parent.getColor() == red {
		uncle := node.getUncle()
		if uncle.getColor() == red {
			// 情况 1：叔叔节点为红。
			node = rbt.fixupRedUncle(node, uncle)
			continue
		}

		if node.parent == node.getGrandparent().left {
			// 情况 2：叔叔黑，父节点在祖父左侧。
			node = rbt.fixupBlackUncleLeftChild(node)
			continue
		}

		// 情况 3：叔叔黑，父节点在祖父右侧。
		node = rbt.fixupBlackUncleRightChild(node)
	}

	// 循环结束后统一保证根节点为黑色。
	rbt.root.setColor(black)
}

// fixupRedUncle 处理“叔叔节点为红”的场景 ( 父叔染黑，祖父染红，继续向上修复 )。
//
//			   b(c)					              r(c)
//			 /	    \					         /     \
//		  r(b)		r(d)					  b(b)	   b(d)
//		 /	\	    /	\		 =>		     /	\	   /   \
//	  b(z)	r(x)  b(p)  b(q)			  b(z)	r(x)  b(p)  b(q)
//
// 1. 将 x 的父节点 b 和叔节点 d 染黑。
func (rbt *RBTree[K, V]) fixupRedUncle(node, uncle *rbNode[K, V]) *rbNode[K, V] {
	grandparent := node.getGrandparent()

	node.parent.setColor(black)
	uncle.setColor(black)
	grandparent.setColor(red)

	return grandparent
}

// fixupBlackUncleLeftChild 处理“叔叔黑 + 父在左” ( 按 LR/LL 做旋转与重染色 )。
//
// 1. 若 x 是父节点的右孩子，先转成 LL 形态 ( 对父节点左旋 )。
//
//			  b(c)				             b(c)
//			 /	  \					        /    \
//		  b(b)	   r(d)					  b(b)	  r(d)
//		 /	\	   /  \	     =>		     /	\	  /   \
//	  b(z)	r(y) b(p) b(q)			   b(z)	r(x) b(p) b(q)
//			   \							       	  /
//		       r(x)								    r(y)
//
// 2. 再围绕祖父做右旋并重染色。
//
//			   b(c)	   	                   b(c)                     r(b)
//			  /	  \	   		              /   \                    /   \
//		  b(b)	   r(d)			      r(b)	   r(d)             b(z)   b(c)
//		 /	\	   /  \	     =>	     /	\	   /  \	     =>            /   \
//	  b(z)	r(x) b(p) b(q)		  b(z)	b(x) b(p) b(q)              b(x)   r(d)
//			/                           /                            /     /  \
//		  r(y)                        r(y)                        r(y)    b(p)  b(q)
//
// 3. 父染黑、祖父染红。
// 4. 祖父右旋，恢复局部平衡与性质。
func (rbt *RBTree[K, V]) fixupBlackUncleLeftChild(node *rbNode[K, V]) *rbNode[K, V] {
	if node == node.parent.right {
		node = node.parent
		rbt.leftRotate(node)
	}

	grandparent := node.getGrandparent()
	node.parent.setColor(black)
	grandparent.setColor(red)
	rbt.rightRotate(grandparent)

	return node.parent
}

// fixupBlackUncleRightChild 处理“叔叔黑 + 父在右” ( 按 RL/RR 做旋转与重染色 )。
//
// 1. 若 x 是父节点的左孩子，先转成 RR 形态 ( 对父节点右旋 )。
//
//			 b(c)						    b(c)
//			/	 \					       /    \
//		  r(b)	  b(d)					  r(b)   b(d)
//		 /	\	  /  \		 =>		     /	\	  /  \
//	  b(p)	b(q) b(z) r(y)			   b(p)	b(q) b(z) r(x)
//			          /                                 \
//		             r(x)                               r(y)
//
// 2. 再围绕祖父做左旋并重染色。
//
//			   b(c)	   	                   b(c)                         r(d)
//			  /	  \	   		              /   \                        /   \
//		  r(b)	   b(d)			      r(b)	   r(d)                 b(c)   b(x)
//		 /	\	   /  \	     =>	     /	\	   /  \	       =>       /  \     \
//	  b(p)	b(q) b(z) r(x)		  b(p)	b(q) b(z) b(x)            r(b) b(z)  r(y)
//			          	\                           \             /  \
//		             	r(y)						r(y)	   b(p)  b(q)
//
// 3. 父染黑、祖父染红。
// 4. 祖父左旋，恢复局部平衡与性质。
func (rbt *RBTree[K, V]) fixupBlackUncleRightChild(node *rbNode[K, V]) *rbNode[K, V] {
	if node == node.parent.left {
		node = node.parent
		rbt.rightRotate(node)
	}

	grandparent := node.getGrandparent()
	node.parent.setColor(black)
	grandparent.setColor(red)
	rbt.leftRotate(grandparent)

	return node
}

// minNode 返回以 node 为根的子树中的最小节点 ( 最左节点 )。
func (rbt *RBTree[K, V]) minNode(node *rbNode[K, V]) *rbNode[K, V] {
	if node == nil || node == rbt.nilNode {
		return nil
	}
	for node.left != nil && node.left != rbt.nilNode {
		node = node.left
	}
	return node
}

// deleteNode 删除给定节点 ( 内部方法 - 调用方已保证节点存在 )。
func (rbt *RBTree[K, V]) deleteNode(node *rbNode[K, V]) {
	// y: 实际被移除的节点 ( 可能是 node 本身，也可能是它的后继 )。
	y := node
	yOriginalColor := y.getColor()
	// x: y 被移除后“顶上来”的节点 ( 可能是哨兵 )。
	var x *rbNode[K, V]

	switch {
	case node.left == rbt.nilNode:
		// 只有右子树或右哨兵。
		x = node.right
		rbt.transplant(node, node.right)
	case node.right == rbt.nilNode:
		// 只有左子树。
		x = node.left
		rbt.transplant(node, node.left)
	default:
		// 双子树：用后继替换。
		y = rbt.minNode(node.right)
		yOriginalColor = y.getColor()
		x = y.right
		if y.parent == node {
			// 后继就是右孩子：只需修正 x 的父指针。
			x.parent = y
		} else {
			// 后继在更深层：先把后继提到它原位置的父节点上。
			rbt.transplant(y, y.right)
			y.right = node.right
			y.right.parent = y
		}

		// 用后继 y 顶替 node 的位置，并继承 node 的颜色。
		rbt.transplant(node, y)
		y.left = node.left
		y.left.parent = y
		y.color = node.color
	}

	if yOriginalColor == black {
		// 仅当移除的是黑节点时，才可能破坏黑高，需要修复。
		rbt.deletionFixup(x)
	}

	if rbt.root == rbt.nilNode {
		// 树被删空时统一恢复为 nil 根，保持外部语义不变。
		rbt.root = nil
	}
}

// transplant 用 v 替换 u 在父节点中的位置 ( 不处理 v 的子树内容 )。
func (rbt *RBTree[K, V]) transplant(u, v *rbNode[K, V]) {
	switch {
	case u.parent == nil:
		rbt.root = v
	case u == u.parent.left:
		u.parent.left = v
	default:
		u.parent.right = v
	}
	if v != nil {
		v.parent = u.parent
	}
}

// deletionFixup 在删除后修复红黑树性质。
//
// 这里使用 CLRS 的经典 4 类情形，对左右分支做镜像处理。
func (rbt *RBTree[K, V]) deletionFixup(node *rbNode[K, V]) {
	for node != rbt.root && node.getColor() == black {
		if node == node.parent.left {
			node = rbt.deletionFixupFromSide(node, true)
		} else {
			node = rbt.deletionFixupFromSide(node, false)
		}
	}

	node.setColor(black)
}

// deletionFixupFromSide 从左侧或右侧修复删除后的红黑树性质。
//
//	nodeIsLeft: 是否从左侧修复。
func (rbt *RBTree[K, V]) deletionFixupFromSide(node *rbNode[K, V], nodeIsLeft bool) *rbNode[K, V] {
	parent := node.parent
	var sibling, nearNephew, farNephew *rbNode[K, V]
	if nodeIsLeft {
		sibling = parent.right
		nearNephew = sibling.left
		farNephew = sibling.right
	} else {
		sibling = parent.left
		nearNephew = sibling.right
		farNephew = sibling.left
	}

	if sibling.getColor() == red {
		// 情况 1：兄弟红，先旋转把问题转成“兄弟黑”再继续 ( 右侧分支为镜像 )。
		sibling.setColor(black)
		parent.setColor(red)
		if nodeIsLeft {
			rbt.leftRotate(parent)
			sibling = parent.right
			nearNephew = sibling.left
			farNephew = sibling.right
		} else {
			rbt.rightRotate(parent)
			sibling = parent.left
			nearNephew = sibling.right
			farNephew = sibling.left
		}
	}

	if sibling.left.getColor() == black && sibling.right.getColor() == black {
		// 情况 2：兄弟黑且两个侄子都黑，向上合并黑高。
		sibling.setColor(red)
		return parent
	}

	if farNephew.getColor() == black {
		// 情况 3：兄弟黑，近侄子红、远侄子黑，先转成情况 4 ( 右侧分支为镜像 )。
		nearNephew.setColor(black)
		sibling.setColor(red)
		if nodeIsLeft {
			rbt.rightRotate(sibling)
			sibling = parent.right
			farNephew = sibling.right
		} else {
			rbt.leftRotate(sibling)
			sibling = parent.left
			farNephew = sibling.left
		}
	}

	// 情况 4：兄弟黑且远侄子红，旋转 + 重染色后结束修复 ( 右侧分支为镜像 )。
	sibling.setColor(parent.getColor())
	parent.setColor(black)
	farNephew.setColor(black)
	if nodeIsLeft {
		rbt.leftRotate(parent)
	} else {
		rbt.rightRotate(parent)
	}
	return rbt.root
}

func (rbt *RBTree[K, V]) findNode(key K) *rbNode[K, V] {
	cmpFunc := rbt.cmp
	node := rbt.root
	for node != nil && node != rbt.nilNode {
		cmp := cmpFunc(key, node.key)
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

// inOrderTraversalBreakable 执行可中断的中序遍历。
// 回调返回 false 时会立刻停止遍历。
func (rbt *RBTree[K, V]) inOrderTraversalBreakable(visitFunc func(node *rbNode[K, V]) bool) {
	stack := make([]*rbNode[K, V], 0, traversalStackCap(rbt.size))
	currNode := rbt.root

	for (currNode != nil && currNode != rbt.nilNode) || len(stack) > 0 {
		for currNode != nil && currNode != rbt.nilNode {
			// 入栈：一路向左走到当前子树最小节点。
			stack = append(stack, currNode)
			currNode = currNode.left
		}

		// 出栈：处理当前节点，再转向右子树。
		currNode = stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if !visitFunc(currNode) {
			break
		}
		currNode = currNode.right
	}
}

// inOrderTraversalAll 执行完整中序遍历 ( 不可中断 )。
// 适合用于 Keys/Vals/Kvs 这类需要收集全量结果的场景。
func (rbt *RBTree[K, V]) inOrderTraversalAll(visitFunc func(node *rbNode[K, V])) {
	stack := make([]*rbNode[K, V], 0, traversalStackCap(rbt.size))

	curr := rbt.root
	for (curr != nil && curr != rbt.nilNode) || len(stack) > 0 {
		for curr != nil && curr != rbt.nilNode {
			stack = append(stack, curr)
			curr = curr.left
		}

		curr = stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		visitFunc(curr)

		curr = curr.right
	}
}

// rbNode 是红黑树节点。
// 考虑内存对齐，颜色字段放在结构体尾部。
type rbNode[K any, V any] struct {
	parent *rbNode[K, V]
	left   *rbNode[K, V]
	right  *rbNode[K, V]
	key    K
	val    V
	color  color
	isNil  bool
}

// newRbNode 创建普通红色节点 ( 不自动挂接哨兵子节点 )。
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

func (rbt *RBTree[K, V]) newNode(key K, val V) *rbNode[K, V] {
	// 所有真实节点的左右孩子默认指向共享哨兵 ( 避免大量 nil 分支判断 )。
	node := newRbNode(key, val)
	node.left = rbt.nilNode
	node.right = rbt.nilNode
	return node
}

func (rbn *rbNode[K, V]) getColor() color {
	if rbn == nil || rbn.isNil {
		return black
	}
	return rbn.color
}

func (rbn *rbNode[K, V]) setColor(color color) {
	if rbn == nil || rbn.isNil {
		return
	}
	rbn.color = color
}

// getGrandparent 返回祖父节点。
func (rbn *rbNode[K, V]) getGrandparent() *rbNode[K, V] {
	if rbn == nil || rbn.parent == nil {
		return nil
	}
	return rbn.parent.parent
}

// getSibling 返回兄弟节点。
func (rbn *rbNode[K, V]) getSibling() *rbNode[K, V] {
	if rbn == nil || rbn.parent == nil {
		return nil
	}
	if rbn == rbn.parent.left {
		return rbn.parent.right
	}
	return rbn.parent.left
}

// getUncle 返回叔叔节点。
func (rbn *rbNode[K, V]) getUncle() *rbNode[K, V] {
	if rbn == nil {
		return nil
	}

	return rbn.parent.getSibling()
}

func traversalStackCap(size int64) int {
	// 红黑树高度满足 h <= 2*log2(n+1) ( 据此估算遍历栈容量 )。
	const minStackCap = 8
	if size <= 0 {
		return 0
	}

	if size <= minStackCap {
		return int(size)
	}

	//nolint:mnd // 计算红黑树高度。
	heightBound := bits.Len64(uint64(size+1)) * 2
	if heightBound < minStackCap {
		return minStackCap
	}
	return heightBound
}
