package chfs

import (
	"strconv"
	"strings"
)

const (
	DirType = iota
	FileType
)

type Object struct {
	// the encoded name... used because this is constant length
	name Name

	objType byte

	executable bool
}

func NewFile(name Name) Object {
	return Object{name,FileType,false}
}

func NewDir(name Name) Object {
	return Object{name,DirType,false}
}

func (o Object) Name() *string {
	return o.name.encoded
}

func (o Object) Type() byte {
	return o.objType
}

func (o Object) repr() string {
	var sb strings.Builder
	sb.WriteString(`"n":"`)
	sb.WriteString(o.name.Encoded())
	sb.WriteString(`","t":`)
	sb.WriteString(strconv.Itoa(int(o.objType)))
	if o.executable {
		sb.WriteString(`,"x":true`)
	}
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
func (b Branch) IsTerminal() bool {
	return b.obj != nil
}

func (b Branch) IsDirectory() bool {
	return b.IsTerminal() && b.obj.objType == DirType
}

func (b Branch) IsFile() bool {
	return b.IsTerminal() && b.obj.objType == FileType
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

func (t Tree) BranchCount() int {
	branchCount := 0
	for _, b := range t.b {
		if b != nil {
			branchCount++
		}
	}
	return branchCount
}

func (t Tree) IsEmpty() bool {
	for _, b := range t.b {
		if b != nil {
			return false
		}
	}
	return true
}

func CopyTree(copy *Tree) *Tree {
	t := new(Tree)
	for i, b := range copy.b {
		if b != nil {
			t.b[i] = CopyBranch(b)
		}
	}
	return t
}

func CopyBranch(copy *Branch) *Branch {
	b := new(Branch)
	b.id = copy.id
	b.obj = copy.obj
	return b
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
