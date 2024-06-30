package chfs

import (
	"fmt"
)

type BasicChFS struct {
	ref Store[string, Ref]
	index Store[Checksum, Index]
	tree Store[Checksum, Tree]
	file Store[Checksum, File]
}

func (fs BasicChFS) GetRef(name string) (Ref, error) {
	return fs.ref.Get(name)
}

func (fs BasicChFS) GetIndex(id Checksum) (Index, error) {
	return fs.index.Get(id)
}

func (fs BasicChFS) GetTree(id Checksum) (Tree, error) {
	return fs.tree.Get(id)
}

func (fs BasicChFS) GetFile(id Checksum) (File, error) {
	return fs.file.Get(id)
}

func (fs BasicChFS) Head() (Ref, error) {
	return fs.ref.Get("HEAD")
}

func (fs BasicChFS) Write(index string, files []FileObj) (Checksum, error) {
	if len(files) != 1 {
		return Checksum{}, fmt.Errorf(`error: basic chfs only supports single file ops`)
	}

	file := files[0]

	root, err := fs.GetRef(index)
	if err != nil {
		return Checksum{}, err
	}
	id := root.id
	return Checksum{}, nil
}

func (fs BasicChFS) Read(index string, paths []Path) ([]FileObj, error) {
	return []FileObj{}, nil
}

func (fs BasicChFS) ListDir(index string, paths []Path) ([][]FileObj, error) {
	return [][]FileObj{}, nil
}
