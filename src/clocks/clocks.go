package clocks

import "time"

// Real

type Clock interface {
	Sleep(time.Duration)
	Since(time.Time) time.Duration
	Now() time.Time
}

func New() Clock {
	return clock{}
}

func (c clock) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (c clock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

func (c clock) Now() time.Time {
	return time.Now()
}

type clock struct{}

// Fake

type FakeClock interface {
	Clock
	Advance(time.Duration)
}

func NewFake(t time.Time) FakeClock {
	return &fakeClock{
		currentTime:   t,
		sleepRequests: []sleepRequest{},
	}
}

func (c *fakeClock) Sleep(d time.Duration) {
	channel := make(chan bool)
	c.sleepRequests = append(c.sleepRequests, sleepRequest{
		goalTime: c.currentTime.Add(d),
		channel:  channel,
	})
	<-channel
}

func (c *fakeClock) Now() time.Time {
	return c.currentTime
}

func (c *fakeClock) Since(t time.Time) time.Duration {
	return c.currentTime.Sub(t)
}

func (c *fakeClock) Advance(d time.Duration) {
	c.currentTime = c.currentTime.Add(d)

	remainingSleepRequests := []sleepRequest{}
	for _, request := range c.sleepRequests {
		if request.goalTime.After(c.currentTime) {
			remainingSleepRequests = append(remainingSleepRequests, request)
		} else {
			request.channel <- true
		}
	}

	c.sleepRequests = remainingSleepRequests
}

type fakeClock struct {
	currentTime   time.Time
	sleepRequests []sleepRequest
}

type sleepRequest struct {
	goalTime time.Time
	channel  chan bool
}
