package fakes

import (
	"sync"
	"time"
)

type FakeClock struct {
	CurrentTime int64
	mutex       sync.Mutex
}

func (fake *FakeClock) Tick() {
	fake.mutex.Lock()
	defer fake.mutex.Unlock()

	fake.CurrentTime++
}

func (fake *FakeClock) TickBy(duration time.Duration) {
	fake.mutex.Lock()
	defer fake.mutex.Unlock()

	fake.CurrentTime += int64(duration.Seconds())
}

func (fake *FakeClock) Now() time.Time {
	fake.mutex.Lock()
	defer fake.mutex.Unlock()

	return time.Unix(fake.CurrentTime, 0)
}

func (fake *FakeClock) Since(t time.Time) time.Duration {
	fake.mutex.Lock()
	defer fake.mutex.Unlock()

	return time.Unix(fake.CurrentTime, 0).Sub(t)
}
