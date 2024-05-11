package vfs

import (
	"fmt"
	"sync"
)

type SimpleVFS struct {
	store Store
	writeMutex sync.Mutex
}

func NewSimpleVFS(store Store) *SimpleVFS {
	vfs := SimpleVFS{}
	vfs.store = store
	return &vfs
}

func (vfs SimpleVFS) Reset() (*Checksum, error) {
	return vfs.setRoot(EmptyBranch())
}

func (vfs SimpleVFS) setRoot(newRoot *Branch) (*Checksum, error) {
	id := newRoot.checksum()
	rootFile := (*File)(id)
	err := vfs.store.Set(&rootKey, rootFile)
	return &rootKey, err
}

func (vfs SimpleVFS) RootKey() (*Checksum, error) {
	f, err := vfs.store.File(&rootKey)
	if err != nil {
		return nil, err
	}
	return f.checksum(), nil
}

func (vfs SimpleVFS) Root() (*Branch, error) {
	rootBranch, err := vfs.RootKey()
	if err != nil {
		return nil, err
	}
	return vfs.store.Branch(rootBranch)
}

func (vfs SimpleVFS) Get(path *Path) (*File, error) {
	curr, err := vfs.RootKey()
	if err != nil {
		return nil, err
	}

	for _, name := range *path {
		for j := 0; j < nameSize; j++ {
			b, err := vfs.store.Branch(curr)
			if err != nil {
				return nil, err
			}
			p := b.next(name.index(j))
			if p == nil {
				// we've reached a terminal node
				return nil, fmt.Errorf("directory does not exist")
			}
			if p.name != nil {
				if name.equals(p.name) {
					break
				} else {
					return nil, fmt.Errorf("directory does not exist")
				}
			}
			curr = &p.id
		}
	}

	return vfs.store.File(curr)
}

func (vfs SimpleVFS) Set(path *Path) (*Checksum, error) {
	vfs.writeMutex.Lock()
	for _, name := range *path {
		for i := 0; i < nameSize; i++ {
		}
	}
	vfs.writeMutex.Unlock()
}
