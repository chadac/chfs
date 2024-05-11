package vfs

import (
	"strings"
)

type Name struct {
	raw string
	encoded Checksum
	IsDir bool
	IsRoot bool
}

type Path []Name

func NewPath(repr string) *Path {
	parts := strings.Split(repr, "/")
	p := (Path)(make([]Name, len(parts)))
	for i, subpath := range parts {
		p[i] = Name{}
		p[i].raw = subpath
		p[i].encoded = encodeString(repr)
		p[i].IsRoot = i == 0
		p[i].IsDir = i < len(parts)-1
	}
	return &p
}

func (p Path) fileName() *Name {
	return &p[len(p)-1]
}

func (n Name) equals(other *Name) bool {
	return n.encoded.equals(&other.encoded)
}

func (n Name) index(index int) byte {
	return n.encoded.index(index)
}