package chfs

type Index struct {
	path string
	parent *Checksum
	id Checksum
}

func (i Index) Key() Checksum {
	return i.id
}
