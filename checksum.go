package chfs

import (
	"crypto/sha256"
	"encoding/hex"
)

const chkSize = 32
type Checksum [32]byte

// convert to a string where each character is an index. it's useful sometimes!
func (c Checksum) Indices() string {
	s := make([]byte, chkSize)
	for i := 0; i < chkSize; i++ {
		s[i] = c.index(i)
	}
	return string(s)
}

func (c Checksum) index(index int) byte {
		return (c[index / 2] >> (4*((index+1) & 1))) & 15
}

func (c Checksum) Equals(that *Checksum) bool {
	for i, b1 := range c {
		if b1 != that[i] {
			return false
		}
	}
	return true
}

func (c Checksum) repr() string {
	return hex.EncodeToString(c[:])
}

func EncodeChecksum(contents string) *Checksum {
	sum := sha256.Sum256([]byte(contents))
	return (*Checksum)(&sum)
}
