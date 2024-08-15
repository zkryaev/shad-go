//go:build !solution

package waitgroup

type WaitGroup struct {
	cnt     chan int
	waiters int
	lock    chan struct{}
}

func New() *WaitGroup {
	wg := &WaitGroup{
		cnt:  make(chan int, 1),
		lock: make(chan struct{}, 1),
	}
	wg.waiters = 0
	wg.cnt <- 0
	return wg
}

func (wg *WaitGroup) Add(delta int) {
	val := <-wg.cnt
	val += delta
	if val < 0 {
		wg.cnt <- val
		panic("negative WaitGroup counter")
	}
	wg.cnt <- val
}

// Do qne decrements the WaitGroup counter by one.
func (wg *WaitGroup) Done() {
	val := <-wg.cnt
	val -= 1
	if val < 0 {
		wg.cnt <- val
		panic("negative WaitGroup counter")
	}
	if val == 0 {
		for wg.waiters > 0 {
			wg.lock <- struct{}{}
			wg.waiters--
		}
		wg.cnt <- val
		return
	}
	wg.cnt <- val
}

// Wait blocks until the WaitGroup counter is zero.
func (wg *WaitGroup) Wait() {
	val := <-wg.cnt
	wg.waiters++
	if val == 0 {
		wg.waiters--
		wg.cnt <- val
		return
	}
	wg.cnt <- val
	<-wg.lock
}
