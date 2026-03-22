package xretry

import "time"

type Strategy interface {
	Next() (time.Duration, bool)
	NextWithRetried(retriedTimes int32) (time.Duration, bool)
	Report(err error) Strategy
}
