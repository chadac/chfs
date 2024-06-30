package chfs

import (
	"fmt"
	"sync"
)

type InMemoryStore[Key comparable] struct {
	mem map[Key]Value[Key]
}

func NewInMemoryStore[Key comparable]() *InMemoryStore[Key] {
	s := InMemoryStore[Key]{}
	s.mem = make(map[Key]Value[Key])
	return &s
}

func (s InMemoryStore[Key]) Get(key Key) (Value[Key], error) {
	value, ok := s.mem[key]
	if ok {
		return value, nil
	} else {
		return nil, fmt.Errorf("key does not exist in store: %v", key)
	}
}

func (s InMemoryStore[Key]) Gets(keys []Key) []Value[Key] {
	values := make([]Value[Key], len(keys))
	var wg sync.WaitGroup
	for i := 0; i <= len(keys); i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			value, err := s.Get(keys[index])
			if err == nil {
				values[index] = value
			}
		}(i)
	}
	wg.Wait()
	return values
}

func (s InMemoryStore[Key]) Put(value Value[Key]) (Key, error) {
	key := value.Key()
	s.mem[key] = value
	return key, nil
}

func (s InMemoryStore[Key]) Puts(values []Value[Key]) []Key {
	keys := make([]Key, len(values))
	var wg sync.WaitGroup
	for i := 0; i <= len(values); i++ {
		wg.Add(1)
		go func(index int) {
			key, err := s.Put(values[index])
			if err == nil {
				keys[index] = key
			}
		}(i)
	}
	return keys
}
