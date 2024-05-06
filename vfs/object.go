package vfs

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// sha256 checksum
type Checksum [32]byte

func (c Checksum) repr() string {
	s := make([]byte, 64)
	for i := 0; i < 64; i++ {
		s[i] = c.index(i)
	}
	return string(s)
}

func (this Checksum) equals(that *Checksum) bool {
	for i, b1 := range this {
		if b1 != that[i] {
			return false
		}
	}
	return true
}

func (c Checksum) index(index int) byte {
	return (c[index / 2] >> (4*((index+1) & 1))) & 15
}

func encodeString(contents string) Checksum {
	sum := sha256.Sum256([]byte(contents))
	return sum
}

type Object interface {
	checksum() Checksum
}

type File string
func (f File) checksum() Checksum {
	return encodeString((string)(f))
}
type branchNode struct {
	id Checksum
	// if not null, points to the next directory in the chain
	name *Name
}
type Branch [16]*branchNode

func (p branchNode) repr() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`{"i":"%x"`, p.id))
	if p.name != nil {
		sb.WriteString(fmt.Sprintf(`,"n":"%x"`, *p.name))
	}
	sb.WriteString("}")
	return sb.String()
}

func (b Branch) repr() string {
	var sb strings.Builder
	sb.WriteString("{")
	for i := 0; i < len(b); i++ {
		if b[i] != nil {
			sb.WriteString(fmt.Sprintf(`"%d":"%s"`, i, b[i].repr()))
			if i != 15 {
				sb.WriteString(",")
			}
		}
	}
	sb.WriteString("}")
	return sb.String()
}

func (b Branch) checksum() Checksum {
	return encodeString(b.repr())
}

func (b Branch) next(index byte) *branchNode {
	return b[index]
}

type Name struct {
	c Checksum
	repr string
}

type Path []Name

func NewPath(repr string) *Path {
	parts := strings.Split(repr, "/")
	p := (Path)(make([]Name, len(parts)))
	for i, subpath := range parts {
		p[i] = Name{}
		p[i].repr = subpath
		p[i].c = encodeString(repr)
	}
	return &p
}

func (p Path) encoded() string {
	var sb strings.Builder
	for _, subpath := range p {
		sb.WriteString(subpath.c.repr())
	}
	return sb.String()
}

func (n Name) equals(other *Name) bool {
	return n.c.equals(&other.c)
}

func (n Name) index(index int) byte {
	return n.c.index(index)
}
