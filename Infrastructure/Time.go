package Infrastructure

import "time"

func NewTime() ITime {
	t := new(Time)
	return t
}

type ITime interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
	Since(t time.Time) time.Duration
}

type Time struct {
}

func (this *Time) Now() time.Time {
	return time.Now()
}

func (this *Time) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func (this *Time) Since(t time.Time) time.Duration {
	return time.Since(t)
}
