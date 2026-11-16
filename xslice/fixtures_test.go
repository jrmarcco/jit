package xslice

// 表驱动用例名在多个测试文件间复用，
// 提取为常量以满足 goconst 约束 ( 出现次数 >= 3 的字面量须为常量 )。
const (
	tcBasic            = "basic"
	tcEmpty            = "empty"
	tcNil              = "nil"
	tcEmptySlice       = "empty slice"
	tcNilSlice         = "nil slice"
	tcNilElems         = "nil elems"
	tcSingleElement    = "single element"
	tcMultipleElements = "multiple elements"
	tcNegativeElements = "negative elements"
	tcMixedElements    = "mixed elements"
	tcContains         = "contains"
	tcNotContains      = "not contains"
	tcNoMatch          = "no match"
	tcFound            = "found"
	tcNotFound         = "not found"
	tcEmptySrc         = "empty src"
	tcEmptyDst         = "empty dst"
	tcIntersection     = "intersection"
	tcNoIntersection   = "no intersection"
)

// personJohn 为 TestReverseOfStruct 使用的样例人名。
const personJohn = "John"
