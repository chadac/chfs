package chfs

type File struct {
	// the checksum of the contents of the file
	id Checksum
}

func (f File) Equals(o *File) bool {
	return f.id.Equals(&o.id)
}

func (f File) Key() Checksum {
	return f.id
}
