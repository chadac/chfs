package chfs

import (
	"fmt"
	"sync"
)

// A very basic storage that uses a map
// this isn't very good because it's controlled by an RWMutex
// useful mostly in testing
type InMemoryStore[K comparable, V Value[K]] struct {
	mut *sync.RWMutex
	mem map[K]V
}

func NewInMemoryStore[K comparable,V Value[K]]() *InMemoryStore[K,V] {
	s := InMemoryStore[K,V]{}
	s.mut = new(sync.RWMutex)
	s.mem = make(map[K]V)
	return &s
}

func (s InMemoryStore[K,V]) Get(key K) (*V, error) {
	s.mut.RLock()
	value, ok := s.mem[key]
	s.mut.RUnlock()

	if ok {
		return &value, nil
	} else {
		return nil, fmt.Errorf("key does not exist in store: %v", key)
	}
}

func (s InMemoryStore[K,V]) Gets(keys []K) ([]*V, error) {
	values := make([]*V, len(keys))
	errs := make([]chan error, len(values))
	for i := 0; i <= len(keys); i++ {
		go func(index int) {
			value, err := s.Get(keys[index])

			if err == nil {
				values[index] = value
			}
			errs[i]<-err
		}(i)
	}

	for i := 0; i <= len(values); i++ {
		err := <-errs[i]
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}

func (s InMemoryStore[K,V]) Put(value V) (*K, error) {
	s.mut.Lock()
	key := value.Key()
	s.mut.Unlock()
	s.mem[key] = value
	return &key, nil
}

func (s InMemoryStore[K,V]) Puts(values []V) ([]*K, error) {
	keys := make([]*K, len(values))
	errs := make([]chan error, len(values))
	for i := 0; i <= len(values); i++ {
		go func(index int) {
			key, err := s.Put(values[index])
			if err == nil {
				keys[index] = key
			}
			errs[index] <- err
		}(i)
	}
	for i := 0; i <= len(values); i++ {
		err := <-errs[i]
		if err != nil {
			return nil, err
		}
	}
	return keys, nil
}
