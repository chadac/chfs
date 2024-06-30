package chfs

import (
	"strconv"
	"strings"
)

const (
	DirType = iota
	FileType
	SymlinkType
)

type Object struct {
	// the encoded name... used because this is constant length
	name string

	objType byte
}

func (o Object) Name() *string {
	return &o.name
}

func (o Object) Type() byte {
	return o.objType
}

func (o Object) repr() string {
	var sb strings.Builder
	sb.WriteString(`"n":"`)
	sb.WriteString(o.name)
	sb.WriteString(`","t":`)
	sb.WriteString(strconv.Itoa(int(o.objType)))
	return sb.String()
}

type Branch struct {
	id Checksum

	// if this is not nil, then it points to a file or directory
	obj *Object
}

func (b Branch) Id() *Checksum {
	return &b.id
}

func (b Branch) repr() string {
	var sb strings.Builder
	sb.WriteString(`{"i": "`)
	sb.WriteString(b.id.repr())
	if b.obj != nil {
		sb.WriteString(",")
		sb.WriteString(b.obj.repr())
	}
	sb.WriteString("}")
	return sb.String()
}

// returns true if this points to a file/directory
func (b Branch) Terminal() bool {
	return b.obj != nil
}

const TreeSize = 16

// Tree is an object used to map escape keys
type Tree struct {
	b [16]*Branch
	id *Checksum
}

func EmptyTree() *Tree {
	t := new(Tree)
	return t
}

// we generate json representations of branches because that's easy
// to consistently hash across  languages
func (t Tree) repr() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, b := range t.b {
		if b == nil {
			sb.WriteString("{}")
		} else {
			sb.WriteString(b.repr())
		}
		if i < len(t.b) - 1 {
			sb.WriteString(",")
		}
	}
	return sb.String()
}

func (t Tree) Key() Checksum {
	if t.id == nil {
		t.id = EncodeChecksum(t.repr())
	}
	return *t.id
}
