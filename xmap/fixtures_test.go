package xmap

// 表驱动用例名在多个测试文件间复用，
// 提取为常量以满足 goconst 约束 ( 出现次数 >= 3 的字面量须为常量 )。
const (
	tcBasic = "basic"
	tcEmpty = "empty"
	tcNil   = "nil"
)
