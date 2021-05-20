package mock_Infrastructure

import (
	"container/list"
	"log"
	"time"
)

func NewFakeITime() *FakeITime {
	fakeTime := new(FakeITime)
	fakeTime.CurrentTime = time.Now()
	return fakeTime
}

type FakeITime struct {
	CurrentTime time.Time
	timers      list.List
}

func newFakeTimer(startTime time.Time, duration time.Duration) *fakeTimer {
	timer := new(fakeTimer)
	timer.channel = make(chan time.Time, 1)
	timer.startTime = startTime
	timer.duration = duration
	return timer
}

type fakeTimer struct {
	channel   chan time.Time
	startTime time.Time
	duration  time.Duration
}

func (this *FakeITime) Now() time.Time {
	return this.CurrentTime
}

func (this *FakeITime) After(d time.Duration) <-chan time.Time {
	timer := newFakeTimer(this.CurrentTime, d)
	this.timers.PushFront(timer)
	return timer.channel
}

func (this *FakeITime) Since(t time.Time) time.Duration {
	return this.CurrentTime.Sub(t)
}

func (this *FakeITime) AdvanceTime(d time.Duration) {
	for {
		nextTimer, nextTimerElement := this.nextTimerToExpireWithinDuration(d)

		if nextTimer == nil {
			log.Println("Found no expired timer")
			this.CurrentTime = this.CurrentTime.Add(d)
			return
		}
		log.Println("Found expired timer")

		elapsedTime := nextTimer.startTime.Add(nextTimer.duration).Sub(this.CurrentTime)
		this.timers.Remove(nextTimerElement)
		this.CurrentTime = this.CurrentTime.Add(elapsedTime)
		d = d - elapsedTime
		nextTimer.channel <- this.CurrentTime
		time.Sleep(time.Millisecond)
	}
}

func (this *FakeITime) nextTimerToExpireWithinDuration(d time.Duration) (*fakeTimer, *list.Element) {
	log.Printf("Finding next timer %d", this.timers.Len())
	var nextTimer *fakeTimer
	var nextTimerElement *list.Element
	for timerElement := this.timers.Front(); timerElement != nil; timerElement = timerElement.Next() {
		log.Println("Looking at first timer")
		timer, ok := timerElement.Value.(*fakeTimer)
		if ok {
			if timer.startTime.Add(timer.duration).Before(this.CurrentTime.Add(d)) {
				if nextTimer == nil {
					nextTimer = timer
					nextTimerElement = timerElement
				} else if timer.startTime.Add(timer.duration).Before(nextTimer.startTime.Add(nextTimer.duration)) {
					nextTimer = timer
					nextTimerElement = timerElement
				}
			}
		}
	}

	return nextTimer, nextTimerElement
}
