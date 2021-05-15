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
	limiter.requestAvailableChannel = make(chan bool, rateLimit)
	limiter.requestSentChannel = make(chan time.Time, rateLimit)
	limiter.endLoopChannel = make(chan bool, 1)

	for i := 0; i < rateLimit; i++ {
		limiter.requestAvailableChannel <- true
	}
	go limiter.run()

	return limiter
}

type rateLimiter struct {
	duration                time.Duration
	requestAvailableChannel chan bool
	requestSentChannel      chan time.Time
	endLoopChannel          chan bool
}

func (this *rateLimiter) PerformInteraction(interaction func()) {
	<-this.requestAvailableChannel
	interaction()
	this.requestSentChannel <- time.Now()
}

func (this *rateLimiter) ShutDown() {
	this.endLoopChannel <- true
}

func (this *rateLimiter) run() {

	for {
		select {
		case <-this.endLoopChannel:
			break
		case requestTime := <-this.requestSentChannel:
			select {
			case <-this.endLoopChannel:
				break
			case <-time.After(this.duration - time.Since(requestTime)):
				this.requestAvailableChannel <- true
			}
		}
	}
}
