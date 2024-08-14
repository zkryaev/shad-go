package rwmutex

type RWMutex struct {
	wlock  chan int
	rlock  chan int
	isBusy bool
}

func New() *RWMutex {
	mutex := &RWMutex{wlock: make(chan int, 1), rlock: make(chan int, 1)}
	mutex.rlock <- 0
	mutex.isBusy = false
	return mutex
}

func (rw *RWMutex) RLock() {
	r := <-rw.rlock
	if !rw.isBusy {
		rw.wlock <- 1
		rw.isBusy = true
	}
	r += 1
	rw.rlock <- r
}

func (rw *RWMutex) RUnlock() {
	r := <-rw.rlock
	if r-1 == 0 {
		rw.isBusy = false
		<-rw.wlock
	}
	r -= 1
	rw.rlock <- r
}

func (rw *RWMutex) Lock() {
	rw.wlock <- 1
}

func (rw *RWMutex) Unlock() {
	<-rw.wlock
}
