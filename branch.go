package vfs

import (
	"fmt"
	"strings"
)

type BranchRef struct {
	id Checksum
	// if not null, used as a name verification to validate that we're
	// on the "proper" terminal node of a tree.
	name *Name
	// if true, this is a terminal node that points to a directory.
	isDir bool

	// logic:
	// isDir & name != nil  => points to root of new directory
	// !isDir & name != nil => points to file
	// name == nil          => internal pointer inside same directory
}
type Branch [16]*BranchRef

func EmptyBranch() *Branch {
	b := new(Branch)
	return b
}

func (p BranchRef) repr() string {
	var sb strings.Builder
	// TODO: faster encoding
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

func (b Branch) key() *Checksum {
	key := EncodeString(b.repr())
	return key
}

func (b Branch) next(index byte) *BranchRef {
	return b[index]
}
