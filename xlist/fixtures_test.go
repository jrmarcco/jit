package xlist

// 跨测试文件复用的表驱动用例名常量，
// 提取为常量以满足 goconst 约束 ( 出现次数 >= 3 的字面量须为常量 )。
// 仅单个测试文件使用的常量就近定义在该文件内。
const (
	tcBasic                      = "basic"
	tcOutOfBounds                = "out of bounds"
	tcBasicOutOfBounds           = "basic out of bounds"
	tcEmptyList                  = "empty list"
	tcInsertToHead               = "insert to head"
	tcInsertToTail               = "insert to tail"
	tcInsertToMiddle             = "insert to middle"
	tcInsertToMiddleOfEmpty      = "insert to middle of empty with size 4"
	tcSetFromHead                = "set from head"
	tcSetFromTail                = "set from tail"
	tcGetFromHead                = "get from head"
	tcGetFromTail                = "get from tail"
	tcDeleteFromHead             = "delete from head"
	tcDeleteFromTail             = "delete from tail"
	tcAppendToNonEmptyList       = "append to non-empty list"
	tcAppendSingleItem           = "append single item"
	tcAppendSingleItemToNonEmpty = "append single item to non-empty list"
)
