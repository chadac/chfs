package vfs

type Index struct {
	name string
	ref Checksum
}

func (i Index) key() Checksum {
	return i.ref
}
