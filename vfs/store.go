package vfs

import (
	"fmt"
)

/**
 * A key-value store that allows for bulk inserts/gets
 **/
type Store interface {
	set(key Checksum, obj Object) error
	Get(id Checksum) (Object, error)
	Gets(ids []Checksum) ([]Object, error)
	Put(obj Object) (Checksum, error)
	Puts(objs []Object) ([]Checksum, error)
}

type InMemoryStore struct {
	mem map[Checksum]Object
}

func NewInMemoryStore() *InMemoryStore {
	s := InMemoryStore{}
	s.mem = make(map[Checksum]Object)
	return &s
}

func (s *InMemoryStore) set(key Checksum, value Object) error {
	s.mem[key] = value
	return nil
}

func (s *InMemoryStore) Get(id Checksum) (Object, error) {
	obj, ok := s.mem[id]
	if ok {
		return obj, nil
	} else {
		return nil, fmt.Errorf("could not find object with key: %x", id)
	}
}

func (s *InMemoryStore) Gets(ids []Checksum) ([]Object, error) {
	objs := make([]Object, len(ids))
	for index, id := range ids {
		obj, ok := s.mem[id]
		if !ok {
			return nil, fmt.Errorf("could not find object with key: %x", id)
		}
		objs[index] = obj
	}
	return objs, nil
}

func (s *InMemoryStore) Put(obj Object) (Checksum, error) {
	id := obj.checksum()
	s.mem[id] = obj
	return id, nil
}

func (s *InMemoryStore) Puts(objs []Object) ([]Checksum, error) {
	ids := make([]Checksum, len(objs))
	for index, obj := range objs {
		id := obj.checksum()
		s.mem[id] = obj
		ids[index] = id
	}
	return ids, nil
}
