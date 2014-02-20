package clocks

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestClocks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Clocks Suite")
}

var _ = Describe("Clocks", func() {
	var clock FakeClock
	var startTime time.Time

	BeforeEach(func() {
		startTime = time.Now()
		clock = NewFake(startTime)
	})

	Describe("sleeping", func() {
		It("blocks until the clock has been advanced by the given duration", func() {
			done := false

			go func() {
				clock.Sleep(5)
				done = true
			}()

			time.Sleep(100)
			Expect(done).To(BeFalse())

			clock.Advance(1)
			time.Sleep(100)
			Expect(done).To(BeFalse())

			clock.Advance(2)
			time.Sleep(100)
			Expect(done).To(BeFalse())

			clock.Advance(2)
			time.Sleep(100)
			Expect(done).To(BeTrue())
		})

		It("works when there are multiple goroutines sleeping", func() {
			done1 := false
			done2 := false

			go func() {
				clock.Sleep(5)
				done1 = true
			}()

			go func() {
				clock.Sleep(8)
				done2 = true
			}()

			time.Sleep(100)

			clock.Advance(2)
			time.Sleep(100)
			Expect(done1).To(BeFalse())
			Expect(done2).To(BeFalse())

			clock.Advance(3)
			time.Sleep(100)
			Expect(done1).To(BeTrue())
			Expect(done2).To(BeFalse())

			clock.Advance(5)
			time.Sleep(100)
			Expect(done2).To(BeTrue())
		})
	})

	Describe("advancing the clock", func() {
		It("updates the value of 'Now'", func() {
			Expect(clock.Now()).To(Equal(startTime))

			clock.Advance(2)
			Expect(clock.Now()).To(Equal(startTime.Add(2)))

			clock.Advance(3)
			Expect(clock.Now()).To(Equal(startTime.Add(5)))
		})
	})
})
