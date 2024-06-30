package chfs

type File struct {
	// the checksum of the file
	id Checksum
}

func (f File) Key() Checksum {
	return f.id
}
