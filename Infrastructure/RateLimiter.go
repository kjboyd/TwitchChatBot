package Infrastructure

import (
	"container/list"
	"fmt"
	"time"
)

/*
Ideally I would mock out the clock so that this class could be tested
*/
type IRateLimiter interface {
	RecordInteraction()
	SleepUntilInteractionAllowed() error
}

func NewRateLimiter(rateLimit int, durationInSeconds int) IRateLimiter {
	limiter := new(rateLimiter)
	limiter.RateLimit = rateLimit
	limiter.DurationInSeconds = durationInSeconds
	limiter.InteractionRecord = list.New()
	return limiter
}

type rateLimiter struct {
	RateLimit         int
	DurationInSeconds int
	InteractionRecord *list.List
}

func (this *rateLimiter) RecordInteraction() {
	this.InteractionRecord.PushFront(time.Now())
}

func (this *rateLimiter) SleepUntilInteractionAllowed() error {
	for {
		this.clearOldInteractionRecords()

		if this.InteractionRecord.Len() >= this.RateLimit {
			callTime, ok := this.InteractionRecord.Back().Value.(time.Time)
			if !ok {
				return fmt.Errorf("Unable to parse api call time.")
			}
			sleepTime := float64(this.DurationInSeconds) - time.Since(callTime).Seconds()
			time.Sleep(time.Duration(sleepTime))
		} else {
			return nil
		}
	}
}

func (this *rateLimiter) clearOldInteractionRecords() {
	for {
		if this.InteractionRecord.Back() == nil {
			return
		}

		callTime, ok := this.InteractionRecord.Back().Value.(time.Time)
		timeSince := time.Since(callTime).Seconds()
		duration := float64(this.DurationInSeconds)
		if ok && timeSince < duration {
			return
		}
		this.InteractionRecord.Remove(this.InteractionRecord.Back())
	}
}
