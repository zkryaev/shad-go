//go:build !solution

package dupcall

import (
	"context"
	"sync"
	"sync/atomic"
)

type Result struct {
	data any
	err  error
}

type Call struct {
	running      bool
	mu           sync.Mutex
	noListenerCh chan struct{}
	cbDoneCh     chan struct{}
	listeners    atomic.Int32
	res          atomic.Pointer[Result]
}

func (o *Call) Do(
	ctx context.Context,
	cb func(context.Context) (interface{}, error),
) (result interface{}, err error) {
	o.listeners.Add(1)
	o.mu.Lock()
	if !o.running {
		o.running = true
		o.noListenerCh = make(chan struct{})
		o.cbDoneCh = make(chan struct{})
		o.mu.Unlock()

		cbCtx, cancelCb := context.WithCancel(context.Background())
		defer cancelCb()

		go func() {
			r := &Result{}
			r.data, r.err = cb(cbCtx)
			o.res.Store(r)
			close(o.cbDoneCh)
		}()

		select {
		case <-o.noListenerCh:
			cancelCb()
			o.mu.Lock()
			o.running = false
			o.mu.Unlock()
			return nil, nil

		case <-o.cbDoneCh:
			o.mu.Lock()
			o.running = false
			o.mu.Unlock()
			r := o.res.Load()
			return r.data, r.err

		case <-ctx.Done():
			o.listeners.Add(-1)
			if o.listeners.Load() == 0 {
				cancelCb()
			}
			return nil, ctx.Err()
		}
	} else {
		o.mu.Unlock()
		select {
		case <-ctx.Done(): // Do canceled
			o.listeners.Add(-1)
			if o.listeners.Load() == 0 {
				close(o.noListenerCh)
			}
		case <-o.cbDoneCh:
			r := o.res.Load()
			return r.data, r.err
		}
	}
	return nil, nil
}

/* v1
func (o *Call) Do(
	ctx context.Context,
	cb func(context.Context) (interface{}, error),
) (result interface{}, err error) {
	res := Result{}
	noListenerCh := make(chan struct{})
	cbDoneCh := make(chan struct{})
	o.mu.Lock()
	if !o.running {
		// if false set to true
		o.running = true
		o.mu.Unlock()
		// to ensure that nobody wait leds to cbStop cb
		cbCtx, cbStop := context.WithCancel(context.Background())
		defer cbStop()
		go func() {
			res.data, res.err = cb(cbCtx)
			close(cbDoneCh)
		}()
		select {
		case <-noListenerCh:
			cbStop()
			return
		case <-cbDoneCh:

			// case check listeners num
			// case wait cb result and close wait channel (sync.Cond to signal them?)
			// cb must run in goroutine
		}
	} else {
		o.mu.Unlock()

		o.listeners.Add(1)
		select {
		case <-ctx.Done(): // Do canceled
			o.listeners.Add(-1)
			if o.listeners.Load() == 0 {
				close(noListenerCh)
			}
		case <-cbDoneCh:

			// somehow track active listeners, and if nobody wait cancel cb
			// case <-call.Result() // cb done
			// unset running
		}
	}
	return res.data, res.err
}
*/
