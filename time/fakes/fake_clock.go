package fakes

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
