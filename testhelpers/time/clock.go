package time

import (
	real_time "time"
)

// TODO: move this into cf/time proper
type Clock interface {
	Tick()
	TickBy(real_time.Duration)
	Now() real_time.Time
}

type FakeClock struct {
	CurrentTime int
}

func (fake *FakeClock) Tick() {
	fake.CurrentTime++
}

func (fake *FakeClock) TickBy(duration real_time.Duration) {
	fake.CurrentTime += duration
}

func (fake *FakeClock) Now() real_time.Time {
	return real_time.Unix(fake.CurrentTime, 0)
}
