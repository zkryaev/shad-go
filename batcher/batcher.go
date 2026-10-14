//go:build !solution

package batcher

import (
	"sync"

	"gitlab.com/slon/shad-go/batcher/slow"
)

type batch struct {
	done   chan struct{}
	result any
}

type Batcher struct {
	mu        *sync.Mutex
	val       *slow.Value
	nextOnce  *sync.Once
	nextBatch *batch
	currBatch *batch

	isRunning bool
}

func NewBatcher(v *slow.Value) *Batcher {
	return &Batcher{
		val:       v,
		nextBatch: &batch{done: make(chan struct{})},
		nextOnce:  &sync.Once{},
		mu:        &sync.Mutex{},
	}
}

func (b *Batcher) Load() any {
	b.mu.Lock()
	prevBatch := b.currBatch
	currentOnce := b.nextOnce
	currentBatch := b.nextBatch
	isRunning := b.isRunning
	b.mu.Unlock()

	if isRunning {
		<-prevBatch.done
	}

	currentOnce.Do(func() {
		b.mu.Lock()
		b.isRunning = true
		b.currBatch = currentBatch
		b.nextBatch = &batch{done: make(chan struct{})}
		b.nextOnce = &sync.Once{}
		b.mu.Unlock()
		res := b.val.Load()

		b.mu.Lock()
		b.isRunning = false
		b.mu.Unlock()

		currentBatch.result = res
		close(currentBatch.done)
	})

	<-currentBatch.done
	return currentBatch.result
}
