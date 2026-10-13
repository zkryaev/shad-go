//go:build !solution

package keylock

import (
	"sync"
)

type KeyLock struct {
	dict  map[string]struct{}
	mu    sync.Mutex
	queue chan struct{}
}

func New() *KeyLock {
	return &KeyLock{
		dict:  make(map[string]struct{}),
		queue: make(chan struct{}),
	}
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	for {
		l.mu.Lock()
		var isAnyLocked bool
		for i := range keys {
			_, ok := l.dict[keys[i]]
			if ok {
				isAnyLocked = true
				break
			}
		}
		if !isAnyLocked {
			for i := range keys {
				l.dict[keys[i]] = struct{}{}
			}
			l.mu.Unlock()
			return false, func() {
				new := make(chan struct{})
				l.mu.Lock()
				for i := range keys {
					delete(l.dict, keys[i])
				}
				old := l.queue
				l.queue = new
				close(old)
				l.mu.Unlock()
			}
		} else {
			queue := l.queue
			l.mu.Unlock()
			select {
			case <-cancel:
				return true, nil
			case <-queue:
			}
		}
	}
}
