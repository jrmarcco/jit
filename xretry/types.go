// Package xretry 提供可插拔退避策略的重试循环。
package xretry

import "time"

// Strategy 决定两次重试之间的等待间隔。
//
// Next 返回下一次重试前的等待时长，以及是否还允许继续重试。
// Strategy 会被 Next 内部计数修改，
// 单个 Strategy 实例不能被多个 Retry 循环并发共享。
type Strategy interface {
	Next() (time.Duration, bool)
}
