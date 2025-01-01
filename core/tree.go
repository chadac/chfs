package "chfs/core"

import (
	"strings"
)

const (
	TreeType = iota,
	FileType
)

type File struct {
	id *Checksum

	// flags
	contentId *Checksum
	executable bool

	// excluded from checksum
	// could be the actual content or a URL to the content
	data string
}

func NewDataFile(data string) *File {
	// create a file where the data in the file is stored in the struct
	return &File{nil,EncodeChecksum(data),false,data}
}

func (f File) contents() string {
	var sb strings.Builder
	sb.WriteString(f.contentId.repr())

	if f.executable {
		sb.WriteString(":x")
	} else {
		sb.WriteString(":!x")
	}

	return sb.String()
}

func (f *File) Id() *Checksum {
	if f.id == nil {
		f.id = EncodeChecksum(f.contents())
	}
	return f.id
}

type Branch struct {
	id *Checksum
	kind uint8
}

func (f *File) AsBranch() *Branch {
	return &Branch{f.Id(), FileType}
}

type Tree struct {
	id *Checksum

	prefix []uint8
	next [17]*Branch
}

func CreateEmptyTree() *Tree {
	return new(Tree)
}

func CreateSingleFileTree(prefix []uint8, file *File) *Tree {
	t := CreateEmptyTree()
	t.prefix = prefix[0:len(prefix)-1]
	t.next[prefix[len(prefix)-1]] = Branch{file.Id(),FileType}
	return t
}

func CreateTwoFileTrees(p1 []uint8, f1 *File, p2 []uint8, f2 *File) []*Tree {
	split := 0
	for ; split < len(p1) & split < len(p2); split++ {
		if p1[split] != p2[split] {
			break
		}
	}
	if p1[split] == p2[split] {
		// if one path is a subset of other, so just overwrite with p2 instead
		return []*Tree{CreateSingleFileTree(p2, f2)}
	} else {
		trees := make([]*Tree, 3)
		t.prefix = p1[0:split]
		t1 := CreateSingleFileTree(p1[split+1:len(p1)], f1)
		t2 := CreateSingleFileTree(p2[split+1:len(p2)], f2)
		t := CreateEmptyTree()
		t.next[p1[split]] = Branch{t1.Id(),TreeType}
		t.next[p2[split]] = Branch{t2.Id(),TreeType}
		return []*Tree{t, t1, t2}
	}
}

func (t Tree) Clone() *Tree {
	newTree := t
	// reset the id for modifications/new checksum computation
	newTree.id = nil
	return &newTree
}

func (t Tree) setBranch(index uint8, branch *Branch) *Tree {
	newTree := t.Clone()
}

// return the number of (non-nil) branches
func (t Tree) NumBranches() uint {
	total := 0
	for i := 0; i < len(t.next); i++ {
		if t.next[i] != nil {
			total++
		}
	}
	return total
}

func (t *Tree) clearNext() {
	for i := 0; i < 17; i++ {
		t.next[i] = nil
	}
}

func (b Branch) appendContents(sb *strings.Builder) {
	if b.kind == TreeType {
		sb.WriteString("t>")
	} else {
		sb.WriteString("f>")
	}
	sb.WriteString(b.id.repr())
}

func (t Tree) contents() string {
	var sb strings.Builder
	for i := 0; i <= 17; i++ {
		if t.next[i] != nil {
			t.next[i].appendContents(sb)
		}
		sb.WriteString(":")
	}
	return sb.String()
}

func (t *Tree) Id() *Checksum {
	if t.id == nil {
		t.id = EncodeChecksum(t.contents())
	}
	return t.id
}

func (t *Tree) AsBranch() *Branch {
	return &Branch{t.Id(),TreeType}
}

func (t Tree) Route(prefix []uint8, store *TreeStore) ([]*Tree, []uint8, uint, error) {
	curr := &t
	route := make([]*Tree, 0)
	branches := make([]uint8, 0)
	i := 0
	var err error
	out:
	for ; i < len(prefix); i++ {
		route = append(route, curr)
		offset := i
		for ; i - offset <= len(curr.prefix) && i < len(prefix); i++ {
			if curr.prefix[i - offset] != prefix[i] {
				break out
			}
		}
		c := prefix[i]
		if curr.next[c] == nil || curr.next[c].kind == FileType {
			break
		}
		branches = append(branches, c)
		curr, err = tree.Get(tree.next[c].id)
		if err != nil {
			return route, branches, c, err
		}
	}
	return route, branches, i, nil
}

func (t1 Tree) sharedPrefix(t2 *Tree) int {
	for i := 0; i < len(t1.prefix); i++ {
		if i >= len(t2.prefix) || t1.prefix[i] != t2.prefix[j] {
			return i
		}
		if t1.prefix[i] != t2.prefix[j] {
			return i
		}
	}
	return len(t1.prefix)
}

// func (t1 Tree) merge(t2 *Tree, store *TreeStore, trees *[]*Tree) (*Tree, error) {
// 	shared := t1.sharedPrefix(t2)
// 	if len(t2.prefix) > len(t1.prefix) && len(t1.prefix) == shared {
// 		// this means that t2 is a subset of t1 that was consolidated
// 	}
// 	if shared == 0 {
// 	} else if shared == 0 {
// 	}
// }

// func (t1 Tree) Merge(t2 *Tree, store *TreeStore) ([]*Tree, error) {
// 	// Combine two trees, favoring t2 during conflicts
// 	// this one's going to be impossible to do iteratively... so gofuncs it is
// }

// func (t1 Tree) Diff(t2 *Tree) ([]*Tree, [][]uint8, error) {
// 	// Compares two trees
// }

func (t Tree) Insert(prefix []uint8, file *File, store *TreeStore) ([]*Tree, error) {
	// inserts using the route method
	// returns a list of the new trees to write, from root to leaves
	// writing should happen in reverse

	// this doesn't use recursion, which may be more helpful in some situations
	route, branches, i, err := t.route(prefix, store)
	if err != nil {
		// read error (corrupt tree), we just return
		return nil, err
	}
	var newTrees []*Tree
	if i < len(prefix)-1 {
		// this means that this path doesn't exist, so we need to insert the rest
		lastTree := route[len(route)-1]
		t1 := CreateSingleFileTree(prefix[i+1:len(prefix)], file)
		t2 := lastTree.clone()
		t2.next[prefix[i]] = t1.AsBranch()
		newTrees = []*Tree{t2, t1}
	} else {
		// TODO: this means that the path does exist and we're overwriting
		if len(route) != len(branches) {
			return nil, fmt.Errorf("Unexpected length received, my logic's borked")
		}
		lastTree := route[len(route)-1]
		t1 := route[len(route)-1].clone()
		t1.next[branches[len(branches)-1]] = file.AsBranch()
		newTrees = []*Tree{t1}
	}

	for j := len(route)-2; j >= 0; j-- {
		newT := route[j].clone()
		newT.next[branches[j]] = route[len(route)-1].AsBranch()
		newTrees = append(newTrees, newT)
	}

	return newTrees, nil
}

func (t Tree) Rm(prefix []uint8, store *TreeStore) ([]*Tree, error) {
	// route-based removal
	route, branches, i, err := t.route(prefix, store)
	if err != nil {
		return nil, err
	} else if i < len(prefix)-1 {
		return nil, fmt.Errorf("file does not exist")
	}
	t1 := route[len(route)-1].clone()
	t1.next[branches[len(branches)-1]] = nil
	newTrees := make([]*Tree, 0)
	for j := len(route)-2; j >= 0; j-- {
		t2 := route[j].clone()
		c := branches[j]
		if len(newTrees) <= 0 && t1.CanPrune() {
			t2.next[c] = nil
		} else {
			t2.next[c] = t1.AsBranch()
			newTrees = append(newTrees, t1)
		}
		t1 = t2
	}

	// note: we always preserve the root tree, so this insert should always happen
	newTrees = append(newTrees, t1)

	return newTrees, nil
}

func (t Tree) Subtree(offset uint) *Tree {
	newTree := t.Clone()
	newTree.prefix = t.prefix[offset:len(t.prefix)]
	return newTree
}

func (t Tree) InsertOld(prefix []uint8, file *File, store *TreeStore) (*Tree, error) {
	// recursive insertion method
	// not really useful so I'm just removing it for now

	// determine if we need a split
	for i := 0; i < len(t.prefix); i++ {
		if i >= len(prefix) {
			// new terminal node! this overwrites presently
			newTree := EmptyTree()
			newTree.prefix = t.prefix[0:i-1]
			newTree.next[prefix[i-1]] = &Branch{file.Id(), FileType}
			return newTree, nil
		}
		else if t.prefix[i] != prefix[i] {
			// split!
			// new shared prefix
			newTree := EmptyTree()
			newTree.prefix = t.prefix[0:i]
			newTree.clearNext()
			newTree.next[t.prefix[i]] = t.Subtree(i)
			newTree.next[prefix[i]] = CreateSingleFileTree(prefix[i:len(prefix)], file)
			return newTree, nil
		}
	}

	// case one: len(t.prefix) == len(prefix)-1
	// this means that we're going to immediately branch out
	if len(t.prefix) == len(prefix) - 1 {
		c := prefix[len(prefix)-1]
		newTree := t.clone()
		newTree.next[c] = &Branch{file.Id(), FileType}
		return newTree
	}

	// case two: we're handing off to where-ever this file is supposed to go
	c := prefix[len(t.prefix)]
	nextTree, err := store.Get(*t.next[c].id)

	if err != nil {
		return nil, err
	}

	newTree := t.clone()
	newTree.next[c] = nextTree.insert(prefix[len(t.prefix):len(prefix)], file, store)
	return newTree, nil
}

// return true if this tree can be pruned
func (t Tree) CanPrune() bool {
	for i := 0; i < 17; i++ {
		if t.next[i] != nil {
			return false
		}
	}
	return true
}

type FileRef struct {
	route []*Tree
	file *File
}

func (f FileRef) fullPath() string {
	parts := make([]string, len(f.route) + 1)
	for i, tree := range f.route {
		parts[i] = tree.label
	}
	parts[len(parts)-1] = f.file.label
	return strings.Join(parts, "")
}
