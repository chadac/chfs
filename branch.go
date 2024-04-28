package vfs

import (
	"fmt"
)

const charSize = 16

type path []directory
type directory struct {
	name *string
	x []uint8
}
type branch [charSize+1]branchNode
type branchNode struct {
	name *string
	ref *id
}

func pathEncode(pathName string) path {
	return path{}
}

func directoryEncode(dirName string) directory {
	return directory{}
}

func initBranch() *branch {
	b := new(branch)
	var i uint8 = 0
	for ; i <= charSize; i++ {
		b[i] = branchNode{nil, nil}
	}
	return b
}

func (b *branch) get(key uint8) (*branchNode, error) {
	if key > charSize {
		return nil, fmt.Errorf("index '%+v' out of range", key)
	}
	return &b[key], nil
}

func (b *branch) set(key uint8, node *branchNode) (*branch, error) {
	if key > charSize {
		return nil, fmt.Errorf("index '%+v' out of range", key)
	}
	newBranch := new(branch)
	var i uint8 = 0
	for ; i <= charSize; i++ {
		if i == key {
			newBranch[i] = *node
		} else {
			newBranch[i] = b[i]
		}
	}
	return newBranch, nil
}

func (b *branch) checksum() *id {
	return nil
}
