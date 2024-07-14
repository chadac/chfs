package chfs

type Ref struct {
	name string
	id Checksum
}

func (r Ref) Key() string {
	return r.name
}

func (r Ref) Id() Checksum {
	return r.id
}
