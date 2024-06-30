package chfs

// import (
// 	"fmt"
// 	"sync"
// )

// // This is an implementation of an SFS that could be used locally
// // without much complication.
// type SimpleSFS struct {
// 	root *Root
// 	store Store
// 	writeMutex *sync.Mutex
// }

// func NewSimpleSFS(store Store) *SimpleSFS {
// 	vfs := SimpleSFS{}
// 	vfs.store = store
// 	vfs.writeMutex = new(sync.Mutex)
// 	return &vfs
// }

// func (vfs SimpleSFS) Reset() (*Checksum, error) {
// 	return vfs.setRoot(EmptyBranch())
// }

// func (vfs SimpleSFS) setRoot(newRoot *Branch) (*Checksum, error) {
// 	newRootKey := newRoot.key()
// 	vfs.root = newRootKey
// 	// rootFile := (*File)(id)
// 	// err := vfs.store.Set(&rootKey, rootFile)
// 	return &newRootKey, nil
// }

// func (vfs SimpleSFS) RootKey() (*Checksum, error) {
// 	return vfs.root, nil
// }

// func (vfs SimpleSFS) Root() (*Branch, error) {
// 	rootBranch, err := vfs.RootKey()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return vfs.store.Branch(rootBranch)
// }

// func (vfs SimpleSFS) Get(path *Path) (*Checksum, error) {
// 	curr, err := vfs.RootKey()
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, name := range *path {
// 		for j := 0; j < nameSize; j++ {
// 			b, err := vfs.store.Branch(curr)
// 			if err != nil {
// 				return nil, err
// 			}
// 			p := b.next(name.index(j))
// 			if p == nil {
// 				// we've reached a terminal node
// 				return nil, fmt.Errorf("directory does not exist")
// 			}
// 			if p.name != nil {
// 				if name.equals(p.name) {
// 					break
// 				} else {
// 					return nil, fmt.Errorf("directory does not exist")
// 				}
// 			}
// 			curr = &p.id
// 		}
// 	}

// 	return vfs.store.File(curr)
// }

// func (vfs SimpleSFS) Write(path *Path, file *Checksum) (*Checksum, error) {
// 	vfs.writeMutex.Lock()
// 	for _, name := range *path {
// 		for i := 0; i < nameSize; i++ {
// 		}
// 	}
// 	vfs.writeMutex.Unlock()
// }
