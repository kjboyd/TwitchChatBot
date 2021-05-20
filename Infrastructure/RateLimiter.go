package Infrastructure

import (
	"time"
)

/*
Ideally I would mock out the clock so that this class could be tested
*/
type IRateLimiter interface {
	PerformInteraction(interaction func())
	ShutDown()
}

func NewRateLimiter(rateLimit int, duration time.Duration, timeObject ITime) IRateLimiter {
	limiter := new(rateLimiter)
	limiter.duration = duration
	limiter.timeObject = timeObject
	limiter.interactionAvailableChannel = make(chan bool, rateLimit)
	limiter.interactionPerformedChannel = make(chan time.Time, rateLimit)
	limiter.endLoopChannel = make(chan bool, 1)

	for i := 0; i < rateLimit; i++ {
		limiter.interactionAvailableChannel <- true
	}
	go limiter.run()

	return limiter
}

type rateLimiter struct {
	duration                    time.Duration
	interactionAvailableChannel chan bool
	interactionPerformedChannel chan time.Time
	endLoopChannel              chan bool
	timeObject                  ITime
}

func (this *rateLimiter) PerformInteraction(interaction func()) {
	<-this.interactionAvailableChannel
	interaction()
	this.interactionPerformedChannel <- this.timeObject.Now()
}

func (this *rateLimiter) ShutDown() {
	this.endLoopChannel <- true
	for len(this.interactionAvailableChannel) != 0 {
		<-this.interactionAvailableChannel
	}
}

func (this *rateLimiter) run() {

	for {
		select {
		case <-this.endLoopChannel:
			return
		case requestTime := <-this.interactionPerformedChannel:
			select {
			case <-this.endLoopChannel:
				return
			case <-this.timeObject.After(this.duration - this.timeObject.Since(requestTime)):
				this.interactionAvailableChannel <- true
			}
		}
	}
}
