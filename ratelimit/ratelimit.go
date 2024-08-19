//go:build !solution

package ratelimit

import (
	"context"
	"errors"
	"time"
)

// Limiter is precise rate limiter with context support.

var ErrStopped = errors.New("limiter stopped")

type Limiter struct {
	NumReqsLimit int
	TimeLimit    time.Duration
	isStoped     bool
	queue        chan struct{}
	lock         chan struct{}
}

func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	l := &Limiter{
		NumReqsLimit: maxCount,
		TimeLimit:    interval,
		queue:        make(chan struct{}, maxCount),
		lock:         make(chan struct{}, 1),
	}
	l.lock <- struct{}{}
	return l
}

func (l *Limiter) Acquire(ctx context.Context) error {
	<-l.lock
	defer func() { l.lock <- struct{}{} }()
	if l.isStoped {
		l.queue = make(chan struct{})
		return ErrStopped
	}
	for {
		select {
		case <-ctx.Done():
			<-l.queue
			return ctx.Err()
		case l.queue <- struct{}{}:
			return nil
		case <-time.After(l.TimeLimit):
			l.queue = make(chan struct{}, l.NumReqsLimit)
		}
	}
}

func (l *Limiter) Stop() {
	l.isStoped = true
}
