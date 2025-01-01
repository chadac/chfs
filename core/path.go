package chfs/core

import (
	"crypto/sha1"
	"strings"
)

type PathName struct {
	simple string
	encoded *string
}

func EncodePathName(raw string) *string {
	sum := sha1.Sum([]byte(raw))
	x := make([]byte, 32)
	for i := 0; i < NameSize; i++ {
		x[i] = sum[i / 2] >> (4*((i+1) & 1)) & 15
	}
	repr := (string)(x[:])
	return &repr
}

func (n PathName) Value() *string {
	if n.encoded == nil {
		n.encoded = EncodePathName(n.simple)
	}
	return n.encoded
}

type Path []uint8

func NewPath(simple string) Path {
}

// type PathTreeNode struct {
// 	prefix []uint8
// 	next []*PathTreeNode
// }

// type PathTreePointer struct {
// 	curr *PathTreeNode
// 	offset uint
// }

// func (p PathTreePointer) next() []PathTreePointer {
// 	if p.offset+1 < len(p.curr.prefix) {
// 		return []PathTreePointer{PathTreePointer{p.curr, p.offset+1}}
// 	} else {
// 		if len(p.curr.next) <= 0 {
// 			return []PathTreePointer{}
// 		}
// 		next := make([]PathTreePointer, len(curr.next))
// 		for i := 0; i < len(next); i++ {
// 			next[i] = PathTreePointer{p.curr.next[i], 0}
// 		}
// 		return next
// 	}
// }

// func (p PathTreePointer) Char() uint8 {
// 	return p.curr.prefix[p.offset]
// }

// func (p PathTreePointer) Rest() []uint8 {
// }
