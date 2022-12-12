package event_loop

import (
	"context"
	"time"
)

type custom_loop_function_t func()

type EventLoop struct {
	auxJobs []custom_loop_function_t
}

func NewEvenlLoop() *EventLoop {
	return &EventLoop{}
}

func (e *EventLoop) AddEvent(f custom_loop_function_t) {
	e.auxJobs = append(e.auxJobs, f)
}

func (e *EventLoop) Schedule(ctx context.Context, p time.Duration, o time.Duration) {
	// Position the first
	next := time.Now().Truncate(p).Add(o)
	if next.Before(time.Now()) {
		next = next.Add(p)
	}
	t := time.NewTimer(time.Until(next))
	for {
		select {
		case <-t.C:
			next = next.Add(p)
			t.Reset(time.Until(next))
			for _, job := range e.auxJobs {
				job()
			}

		case <-ctx.Done():
			t.Stop()
			return
		}
	}
}
