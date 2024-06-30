package chfs


type SimpleChFS struct {
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

func (fs SimpleChFS) Write(index string, files []FileObj) (Checksum, error) {
	return Checksum{}, nil
}

func (fs SimpleChFS) Read(index string, paths []Path) ([]FileObj, error) {
	return []FileObj{}, nil
}

func (fs SimpleChFS) ListDir(index string, paths []Path) ([][]FileObj, error) {
	return [][]FileObj{}, nil
}
