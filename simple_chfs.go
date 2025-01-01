package chfs

import (
	// "fmt"
	"sync"
)

/// SimpleChFS: A basic implementation of ChFS for single writes
type SimpleChFS struct {
	writeMu *sync.Mutex

	ref Store[string, Ref]
	index Store[Checksum, Index]
	tree Store[Checksum, Tree]
	file Store[Checksum, File]
}

func NewSimpleChFS() *SimpleChFS {
	fs := SimpleChFS{}
	fs.writeMu = new(sync.Mutex)
	fs.ref = NewInMemoryStore[string, Ref]()
	fs.index = NewInMemoryStore[Checksum, Index]()
	fs.tree = NewInMemoryStore[Checksum, Tree]()
	fs.file = NewInMemoryStore[Checksum, File]()
	return &fs
}

func (fs SimpleChFS) GetRef(name string) (*Ref, error) {
	return fs.ref.Get(name)
}

func (fs SimpleChFS) GetIndex(id Checksum) (*Index, error) {
	return fs.index.Get(id)
}

func (fs SimpleChFS) GetTree(id Checksum) (*Tree, error) {
	return fs.tree.Get(id)
}

func (fs SimpleChFS) GetFile(id Checksum) (*File, error) {
	return fs.file.Get(id)
}

func (fs SimpleChFS) Head() (*Ref, error) {
	return fs.ref.Get("HEAD")
}

func (fs *SimpleChFS) Write(refName string, actions []PathObject) (*Checksum, error) {
	paths := make([]Path, len(actions))
	branches := make([]*Branch, len(actions))
	for i, a := range actions {
		paths[i] = a.Path()
		branches[i] = a.Branch()
	}

	fs.writeMu.Lock()
	ref, err := fs.ref.Get(refName)
	if err != nil {
		return nil, err
	}
	subtree := NewSubTree(paths, NewWriter)

	// now for our three steps
	// TODO: make this fast!
	ReadTree(subtree, fs.tree, ref.id)
	Plan(subtree, branches)

	subtree.Print()

	newRoot := WriteTree(subtree, fs.tree)

	// fmt.Printf("%s\n", newRoot.repr())
	// update our tree root
	fs.ref.Put(Ref{refName,newRoot.Key()})

	fs.writeMu.Unlock()

	key := newRoot.Key()
	return &key, nil
}

func (fs SimpleChFS) Read(refName string, paths []Path) ([]PathObject, error) {
	ref, err := fs.ref.Get(refName)
	if err != nil {
		return []PathObject{}, err
	}

	subtree := NewSubTree(paths, NewReader)

	pathObjects := make([]PathObject, len(paths))
	ReadTree(subtree, fs.tree, ref.id)
	for i, leaf := range subtree.leafs {
		pathObjects[i] = NewPathObject(leaf.Dir(), leaf.extra.CurrFile())
	}

	return pathObjects, nil
}

func (fs *SimpleChFS) ListDir(refName string, path Path) ([]PathObject, error) {
	ref, err := fs.Head()
	if err != nil {
		return nil, err
	}
	files, err := ListDir(fs.tree, ref.id, true)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (fs *SimpleChFS) Tree() error {
	return nil
}
