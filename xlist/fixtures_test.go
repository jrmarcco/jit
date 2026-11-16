package xlist

// 表驱动用例名在多个测试文件间复用，
// 提取为常量以满足 goconst 约束 ( 出现次数 >= 3 的字面量须为常量 )。
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
