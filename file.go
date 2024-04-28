package vfs

type file struct {
	contents *string
}

func initFile() *file {
	f := new(file)
	f.contents = nil
	return f
}

func (f *file) checksum() *id {
	return nil
}
