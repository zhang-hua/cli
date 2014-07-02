package fakes

import "time"

type FakeClock struct {
	CurrentTime int64
}

func (fake *FakeClock) Tick() {
	fake.CurrentTime++
}

func (fake *FakeClock) TickBy(duration time.Duration) {
	fake.CurrentTime += int64(duration.Seconds())
}

func (fake *FakeClock) Now() time.Time {
	return time.Unix(fake.CurrentTime, 0)
}
