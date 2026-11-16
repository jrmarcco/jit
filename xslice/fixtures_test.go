package xslice

// 跨测试文件复用的表驱动用例名常量，
// 提取为常量以满足 goconst 约束 ( 出现次数 >= 3 的字面量须为常量 )。
// 仅单个测试文件使用的常量就近定义在该文件内。
const (
	tcBasic            = "basic"
	tcEmpty            = "empty"
	tcEmptySlice       = "empty slice"
	tcNilSlice         = "nil slice"
	tcSingleElement    = "single element"
	tcMultipleElements = "multiple elements"
	tcNoMatch          = "no match"
	tcEmptySrc         = "empty src"
	tcEmptyDst         = "empty dst"
	tcIntersection     = "intersection"
	tcNoIntersection   = "no intersection"
)
