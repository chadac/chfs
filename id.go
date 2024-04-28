package vfs

import (
	"crypto/sha256"
)

const chkSize = sha256.Size
type id [chkSize]byte

func mkRootKey() id {
	rootKey := id{}
	for i := 0; i <= chkSize; i++ {
		rootKey[i] = 0x0
	}
	return rootKey
}
var rootKey = mkRootKey()

func (this *id) equals(other *id) bool {
	for i := 0; i < chkSize; i++ {
		if this[i] != other[i] {
			return false
		}
	}
	return true
}
