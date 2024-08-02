//go:build !solution

package lrucache

import (
	"container/list"
	"fmt"
)

type Content struct {
	key   int
	value int
}

type LRUcache struct {
	keys   map[int]*list.Element
	values list.List
	cap    int
	len    int
}

func (lru *LRUcache) printList() {
	fmt.Println("printing list....")
	for l := lru.values.Front(); l != nil; l = l.Next() {
		fmt.Println(l, "______")
		fmt.Println("n:", l.Next())
		fmt.Println("p:", l.Prev())
	}
}

func (lru *LRUcache) sanitize() bool {
	if lru.len <= 0 {
		return false
	}
	LRUelement := lru.values.Remove(lru.values.Back())
	LRUKey := LRUelement.(Content).key
	delete(lru.keys, LRUKey)
	lru.len -= 1
	return true
}

func (lru *LRUcache) Get(key int) (int, bool) {
	if lru.len <= 0 {
		return 0, false
	}
	GotElement, ok := lru.keys[key]
	if !ok { // key doesnt exist
		return 0, false
	}
	GotItem := (((*GotElement).Value).(Content))
	lru.values.MoveToFront(GotElement)
	return GotItem.value, true
}

func (lru *LRUcache) Set(key, value int) {
	if lru.len == lru.cap {
		ok := lru.sanitize()
		if !ok {
			return
		}
	}
	NewItem := Content{key: key, value: value}
	if _, ok := lru.keys[key]; !ok {
		lru.keys[key] = lru.values.PushFront(NewItem)
		lru.len++
	} else {
		lru.values.Remove(lru.keys[key])
		lru.keys[key] = lru.values.PushFront(NewItem)
	}
}

func (lru *LRUcache) Range(f func(key, value int) bool) {
	if lru.len <= 0 {
		return
	}
	seeker := lru.values.Back()
	for i := 0; i < lru.values.Len(); i++ {
		if seeker != nil && !f(seeker.Value.(Content).key, seeker.Value.(Content).value) {
			break
		}
		seeker = seeker.Prev()
	}

}

func (lru *LRUcache) Clear() {
	if lru.len <= 0 {
		return
	}
	lru.keys = make(map[int]*list.Element)
	lru.values.Init()
	lru.len = 0
}

func New(cap int) Cache {
	cache := &LRUcache{
		keys:   make(map[int]*list.Element, cap),
		values: *list.New(),
		cap:    cap,
		len:    0,
	}
	cache.values.Init()
	return cache
}
