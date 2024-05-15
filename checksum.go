package vfs

import (
	"crypto/sha256"
)

// sha256 checksum
type Checksum [32]byte

const branchSize = 64;
const nameSize = 64;

// convert to a string where each character is an index. it's useful sometimes!
func (c Checksum) ToString() string {
	s := make([]byte, nameSize)
	for i := 0; i < nameSize; i++ {
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

func EncodeString(contents string) *Checksum {
	sum := sha256.Sum256([]byte(contents))
	return (*Checksum)(&sum)
}
