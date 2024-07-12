package chfs

import (
	"sync"
)

/// SimpleChFS: A basic implementation of ChFS for single writes
type SimpleChFS struct {
	writeMu sync.Mutex

	ref Store[string, Ref]
	index Store[Checksum, Index]
	tree Store[Checksum, Tree]
	file Store[Checksum, File]
}

func (fs SimpleChFS) GetRef(name string) (Ref, error) {
	return fs.ref.Get(name)
}

func (fs SimpleChFS) GetIndex(id Checksum) (Index, error) {
	return fs.index.Get(id)
}

func (fs SimpleChFS) GetTree(id Checksum) (Tree, error) {
	return fs.tree.Get(id)
}

func (fs SimpleChFS) GetFile(id Checksum) (File, error) {
	return fs.file.Get(id)
}

func (fs SimpleChFS) Head() (Ref, error) {
	return fs.ref.Get("HEAD")
}

func (fs SimpleChFS) Write(index string, actions []FileObj) (Checksum, error) {
	paths := make([]Path, len(actions))
	files := make([]File, len(actions))
	for i, a := range actions {
		paths[i] = a.path
		files[i] = a.file
	}

	fs.writeMu.Lock()
	ref, _ := fs.ref.Get(index)
	subtree := NewSubTree(paths, false)
	plan := NewPlan(subtree)
	// now for our three steps
	// TODO: this may be more involved in the future
	plan.Read(ref.Id(), fs.tree)
	plan.Update(files)
	plan.Write(fs.tree)

	fs.writeMu.Unlock()
	return Checksum{}, nil
}

func (fs SimpleChFS) Read(index string, paths []Path) ([]FileObj, error) {
	// ref, _ := fs.ref.Get(index)

	// subtree := NewSubTree(paths, false)

	return []FileObj{}, nil
}

func (fs SimpleChFS) ListDir(index string, paths []Path) ([][]FileObj, error) {
	return [][]FileObj{}, nil
}
