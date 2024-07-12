package chfs

import (
	"crypto/sha1"
	"strings"
)

const NameSize = 40

type Name struct {
	raw string
	encoded *string
}

func encodeName(raw string) *string {
	sum := sha1.Sum([]byte(raw))
	x := make([]byte, NameSize)
	for i := 0; i < NameSize; i++ {
		x[i] = sum[i / 2] >> (4*((i+1) & 1)) & 15
	}
	repr := (string)(x[:])
	return &repr
}

func NewName(raw string) *Name {
	n := Name{raw,nil}
	return &n
}

func (n Name) Encoded() string {
	if n.encoded == nil {
		n.encoded = encodeName(n.raw)
	}
	return *n.encoded
}

func (n Name) Index(index int) byte {
	return n.Encoded()[index]
}

type Path []*Name

func NewPath(repr string) Path {
	// remove leading forward slash
	if repr[0] == '/' {
		repr = repr[1:]
	}
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
