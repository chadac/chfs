package vfs

import (
	"sync"
)

type WriteCache interface {
	Id() int64
	Claim(path string) int64
	Acquire(subpath string, id int64) *Branch
	Release(subpath string, b *Branch)
}

type cacheEntry struct {
	id int64
	c chan Branch
}
type cachePriorityQueue []*cacheEntry

func (pq *cachePriorityQueue) Push(x *cacheEntry) {
	*pq = append(*pq, x)
}

func (pq *cachePriorityQueue) Pop() *cacheEntry {
	old := *pq
	n := len(old)
	var next *cacheEntry = nil
	if n > 2 {
		next = old[n-2]
	}
	*pq = old[0 : n-1]
	return next
}

func (pq cachePriorityQueue) GetById(id int64) (int, *cacheEntry) {
	for index, entry := range pq {
		if entry.id == id {
			return index, entry
		}
	}
	return -1, nil
}

type LocalWriteCache struct {
	mu sync.Mutex
	id int64
	m map[string]*cachePriorityQueue
}

func NewLocalWriteCache() *LocalWriteCache {
	c := LocalWriteCache{}
	c.id = 0
	c.m = make(map[string]*cachePriorityQueue)
	return &c
}

// Generate ID to ensure that all changes occur in order
func (c LocalWriteCache) newId() int64 {
	id := c.id
	c.id += 1
	return id
}

func (c LocalWriteCache) set(key string, id int64) {
	// condition: set is always called with the highest id
	pq, ok := c.m[key]
	if !ok {
		e := cacheEntry{id, nil}
		pq = &cachePriorityQueue{&e}
		c.m[key] = pq
	} else {
		// if an entry already exists, then we need to set up ourselves to
		// wait in line until we're up
		e := cacheEntry{id, make(chan Branch)}
		pq.Push(&e)
	}
}

func (c LocalWriteCache) Claim(path string) int64 {
	c.mu.Lock()
	id := c.newId()
	for i := 0; i < len(path); i++ {
		c.set(path[:i], id)
	}
	c.mu.Unlock()
	return id
}

func (c LocalWriteCache) Acquire(subpath string, id int64) *Branch {
	var b *Branch = nil
	var ch chan Branch = nil

	c.mu.Lock()
	pq := c.m[subpath]
	c.mu.Unlock()

	i, e := pq.GetById(id)
	if i == -1 {
		// this should never happen! error
	} else if i > 0 {
		// means we're waiting on someone else
		ch = e.c
	}

	// wait until we are next in line
	if ch != nil {
		newBranch := <-ch
		b = &newBranch
	}

	return b
}

func (c LocalWriteCache) Release(subpath string, b *Branch) {
	c.mu.Lock()
	pq := c.m[subpath]
	// the first entry in this queue is us
	// the next entry is whomever is waiting on us
	next := pq.Pop()

	// if the queue is empty then we can go ahead and delete the entry
	if len(*pq) == 0 {
		delete(c.m, subpath)
	}
	c.mu.Unlock()

	// notify the next in line that they're ready to go
	if next != nil {
		next.c <- *b
	}
}
