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

func (n Name) String() string {
	return n.raw
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

func EmptyPath() Path {
	return make([]*Name, 0)
}

func NewPath(repr string) Path {
	// remove leading forward slash
	if repr[0] == '/' {
		repr = repr[1:]
	}
	if len(repr) == 0 {
		return EmptyPath()
	}
	parts := strings.Split(repr, "/")
	p := (Path)(make([]*Name, len(parts)))
	for i, subpath := range parts {
		p[i] = NewName(subpath)
	}
	return p
}

func (p Path) String() string {
	if len(p) > 0 {
		s := make([]string, len(p))
		for i, sp := range p {
			s[i] = sp.String()
		}
		return "/" + strings.Join(s, "/")
	} else {
		return ""
	}
}

func (p Path) Base() *Name {
	if len(p) > 0 {
		return p[len(p)-1]
	} else {
		return nil
	}
}

func (p Path) Append(n Name) Path {
	return append(p[:], &n)
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
