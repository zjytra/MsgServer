/*
创建时间: 2021/8/26 21:28
作者: zjy
功能介绍:

*/

package timingwheel

import "time"

type EveryScheduler struct {
	Interval time.Duration
}

func (s *EveryScheduler) Next(prev time.Time) time.Time {
	return prev.Add(s.Interval)
}

