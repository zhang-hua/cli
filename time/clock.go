package time

import (
	real_time "time"
)

type Clock interface {
	Now() real_time.Time
}

type realClock struct{}

func NewClock() Clock {
	return realClock{}
}

func (clock realClock) Now() real_time.Time {
	return real_time.Now()
}
