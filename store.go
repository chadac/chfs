package vfs

// type Value[Key comparable] interface {
// 	key() *Key
// }

// type CalculateKeyFunc func(comparable, Value)

// /**
//  * A key-value store that allows for bulk inserts/gets
//  **/
// type Store[Key comparable] interface {
// 	Get(key *Key) (Value[Key], error)
// 	Gets(keys []*Key) ([]Value[Key], error)
// 	Put(value Value[Key]) (Key, error)
// 	Puts(values []Value[Key]) ([]Key, error)
// }

// func StoreGetBranch(store Store, id *Checksum) (*Branch, error) {
// 	obj, err := store.Get(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	branch, ok := obj.(*Branch)
// 	if !ok {
// 		return nil, fmt.Errorf("object at location '%x' is not a branch", id)
// 	}
// 	return branch, nil
// }

// type InMemoryStore[Key comparable] struct {
// 	mem map[Key]Value
// }

// func NewInMemoryStore() *InMemoryStore {
// 	s := InMemoryStore{}
// 	s.mem = make(map[Checksum]Object)
// 	return &s
// }

// func (s *InMemoryStore) Set(key *Checksum, value Object) error {
// 	s.mem[*key] = value
// 	return nil
// }

// func (s *InMemoryStore) Get(id *Checksum) (Object, error) {
// 	obj, ok := s.mem[*id]
// 	if ok {
// 		return obj, nil
// 	} else {
// 		return nil, fmt.Errorf("could not find object with key: %x", id)
// 	}
// }

// func (s InMemoryStore) Branch(id *Checksum) (*Branch, error) {
// 	return StoreGetBranch(&s, id)
// }

// func (s *InMemoryStore) Gets(ids []*Checksum) ([]Object, error) {
// 	objs := make([]Object, len(ids))
// 	for index, id := range ids {
// 		obj, ok := s.mem[*id]
// 		if !ok {
// 			return nil, fmt.Errorf("could not find object with key: %x", id)
// 		}
// 		objs[index] = obj
// 	}
// 	return objs, nil
// }

// func (s *InMemoryStore) Put(obj Object) (*Checksum, error) {
// 	id := obj.key()
// 	s.mem[*id] = obj
// 	return id, nil
// }

// func (s *InMemoryStore) Puts(objs []Object) ([]*Checksum, error) {
// 	ids := make([]*Checksum, len(objs))
// 	for index, obj := range objs {
// 		id := obj.key()
// 		s.mem[*id] = obj
// 		ids[index] = id
// 	}
// 	return ids, nil
// }
