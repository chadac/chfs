package vfs

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// sha256 checksum
type Checksum [32]byte

const branchSize = 64;
const nameSize = 64;

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
	checksum() *Checksum
	file() *File
	branch() *Branch
}

type File Checksum
func FileFromString(contents string) *File {
	chk := encodeString(contents)
	return (*File)(&chk)
}
func (f File) checksum() *Checksum {
	return (*Checksum)(&f)
}
func (f File) file() *File {
	return &f
}
func (f File) branch() *Branch {
	return nil
}

type branchNode struct {
	id Checksum
	// if not null, this is a name verification mechanism
	name *Name
	// if true, mark as directory pointer
	isDir bool
}
type Branch [16]*branchNode

func EmptyBranch() *Branch {
	b := Branch{}
	for i := 0; i < 16; i++ {
		b[i] = nil
	}
	return &b
}

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

func (b Branch) checksum() *Checksum {
	chk := encodeString(b.repr())
	return &chk
}

func (b Branch) file() *File {
	return nil
}

func (b Branch) branch() *Branch {
	return &b
}

func (b Branch) next(index byte) *branchNode {
	return b[index]
}
