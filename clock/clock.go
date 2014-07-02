package clock

import "time"

type Clock interface {
	Now() time.Time
}

type realClock struct{}

func NewClock() Clock {
	return realClock{}
}

func (clock realClock) Now() time.Time {
	return time.Now()
}
