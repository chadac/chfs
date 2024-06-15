package vfs

import (
	"strings"
)

const NameSize = 32

type Name struct {
	raw string
	encoded *string
}

func encodeName(raw string) *string {
	return nil
}

func NewName(raw string) *Name {
	n := Name{raw,nil}
	return &n
}

func (n Name) Index(index int) byte {
	if n.encoded == nil {
		n.encoded = encodeName(n.raw)
	}
	return (*n.encoded)[index]
}

type Path []*Name

func NewPath(repr string) Path {
	parts := strings.Split(repr, "/")
	p := (Path)(make([]*Name, len(parts)))
	for i, subpath := range parts {
		p[i] = NewName(subpath)
	}
	return p
}

func (p Path) Base() *Name {
	return p[len(p)-1]
}

// type NameIndex struct {
// 	n *Name
// 	i byte
// }

// // Returns path representation
// func PathsToTree(paths []*Path) (NameIndex, NameIndex)[] {
// 	var edges := make((NameIndex, NameIndex)[], len(paths))
// 	pathBuf
// 	for i := 0; true; i++ {
// 	}
// 	return edges
// }
