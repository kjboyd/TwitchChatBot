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

func NewRateLimiter(rateLimit int, duration time.Duration) IRateLimiter {
	limiter := new(rateLimiter)
	limiter.duration = duration
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
}

func (this *rateLimiter) PerformInteraction(interaction func()) {
	<-this.interactionAvailableChannel
	interaction()
	this.interactionPerformedChannel <- time.Now()
}

func (this *rateLimiter) ShutDown() {
	this.endLoopChannel <- true
}

func (this *rateLimiter) run() {

	for {
		select {
		case <-this.endLoopChannel:
			break
		case requestTime := <-this.interactionPerformedChannel:
			select {
			case <-this.endLoopChannel:
				break
			case <-time.After(this.duration - time.Since(requestTime)):
				this.interactionAvailableChannel <- true
			}
		}
	}
}
