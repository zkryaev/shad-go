//go:build !solution

package cond

type Locker interface {
	Lock()
	Unlock()
}

type Cond struct {
	L        Locker
	queue    chan struct{}
	abonents chan int
}

func New(l Locker) *Cond {
	cond := &Cond{
		L:        l,
		queue:    make(chan struct{}, 1),
		abonents: make(chan int, 1),
	}
	cond.abonents <- 0
	cond.queue <- struct{}{}
	return cond
}

func (c *Cond) Wait() {
	n := <-c.abonents
	c.L.Unlock()
	if c.L == nil {
		c.abonents <- n
		panic("L didn't hold")
	}
	n += 1
	c.abonents <- n
	c.queue <- struct{}{}
	c.L.Lock()
}

func (c *Cond) Signal() {
	n := <-c.abonents
	if n == 0 {
		c.abonents <- n
		return
	}
	n -= 1
	c.abonents <- n
	<-c.queue
}

func (c *Cond) Broadcast() {
	n := <-c.abonents
	for n != 0 {
		<-c.queue
		n--
	}
	c.abonents <- n
}
