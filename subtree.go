package chfs

import (
	"fmt"
)

// TODO: Swap pnode to use this
type DirInfo[T any] struct {
	root *PNode[T]
	prevDir *PNode[T]
	nextDir *PNode[T]
	prevTail *PNode[T]
}

type PNode[T any] struct {
	subtree *SubTree[T]

	dirPath Path
	name *Name

	pathIndex int
	nameIndex int

	char byte

	prev *PNode[T]
	next []*PNode[T]

	// points to the root of the current directory
	dir *PNode[T]
	// points to the head of the previous directory (if it exists)
	prevDir *PNode[T]
	// points to the head of the next directory (if it exists)
	nextDir *PNode[T]
	// points to the last node of the previous directory
	prevTail *PNode[T]

	extra T
}

func (n PNode[T]) String() string {
	if n.IsRoot() {
		return "{subtree:/[0]}"
	} else {
		return fmt.Sprintf("{subtree:%s/%s[%d].%d}", n.Dir(), n.Name(), n.nameIndex, n.char)
	}
}

// a PathGroup is a tree representing the combined representations of
// several paths
type SubTree[T any] struct {
	root *PNode[T]
	paths []Path
	leafs []*PNode[T]

	newExtra func() T
}

func NewSubTree[T any](paths []Path, newExtra func() T) *SubTree[T] {
	t := SubTree[T]{}

	t.newExtra = newExtra
	t.paths = make([]Path, 0)
	t.leafs = make([]*PNode[T], 0)

	// init root node
	root := PNode[T]{}
	root.subtree = &t
	root.pathIndex = 0
	root.nameIndex = -1
	root.next = make([]*PNode[T], 0)
	root.extra = t.newExtra()
	t.root = &root

	for _, p := range paths {
		t.AddPath(p)
	}

	return &t
}

func (t *SubTree[T]) AddPath(p Path) {
	newNodes := t.root.addPath(p)
	if newNodes != nil {
		leaf := &newNodes[len(newNodes)-1]
		t.leafs = append(t.leafs, leaf)
		t.paths = append(t.paths, p)
	}
}

func (n PNode[T]) nextIndices() (int, int) {
	nextPathIndex := n.pathIndex
	nextNameIndex := n.nameIndex + 1
	if nextNameIndex == NameSize {
		// terminal node, move to the next
		nextPathIndex += 1
		nextNameIndex = 0
	}
	return nextPathIndex, nextNameIndex
}

func (n *PNode[T]) addPath(p Path) []PNode[T] {
	nextPathIndex, nextNameIndex := n.nextIndices()
	if nextPathIndex >= len(p) {
		// we apparently already added this path, so skip
		return nil
	}

	pchar := p[nextPathIndex].Index(nextNameIndex)
	for _, next := range n.next {
		if pchar == next.char {
			return next.addPath(p)
		}
	}

	newNodes := n.createNodeChain(p, nextPathIndex, nextNameIndex)
	n.next = append(n.next, &newNodes[0])
	return newNodes
}

func (prev *PNode[T]) createNodeChain(p Path, pathIndex int, nameIndex int) []PNode[T] {
	t := prev.subtree

	// generates and links together a chain of nodes
	numNodes := (len(p)-pathIndex-1)*NameSize + (NameSize-nameIndex)
	nodes := make([]PNode[T], numNodes)

	nodes[0].prev = prev
	startIndex := pathIndex * NameSize + nameIndex
	for i := 0; i < len(nodes); i++ {
		n := &nodes[i]

		n.subtree = t

		// initialize extra values
		extraValue := t.newExtra()
		n.extra = extraValue

		// set pathIndex, nameIndex appropriately
		totalIndex := pathIndex * NameSize + nameIndex + i
		n.pathIndex = totalIndex / NameSize
		n.nameIndex = totalIndex % NameSize

		n.char = p[n.pathIndex].Index(n.nameIndex)
		n.dirPath = p[:n.pathIndex]
		n.name = p[n.pathIndex]

		dirIndex := n.pathIndex * NameSize - startIndex
		if dirIndex >= 0 {
			n.dir = &nodes[dirIndex]
		} else if n.pathIndex == prev.pathIndex {
			n.dir = prev.dir
		} else {
			panic("reached an unanticipated edge case! test better")
		}

		prevDirIndex := (n.pathIndex-1) * NameSize - startIndex
		if prevDirIndex >= 0 {
			n.prevDir = &nodes[prevDirIndex]
		} else if n.pathIndex == prev.pathIndex {
			// case 1: we share previous directories
			n.prevDir = prev.prevDir
		} else if n.pathIndex == prev.pathIndex+1 {
			// case 2: we point to prev's root directory
			n.prevDir = prev.dir
		} else {
			panic("reached an unanticipated edge case! test better")
		}

		prevTailIndex := prevDirIndex + NameSize - 1
		if prevTailIndex >= 0 {
			n.prevTail = &nodes[prevTailIndex]
		} else if prevTailIndex == -1 {
			n.prevTail = prev
		} else if n.pathIndex == prev.pathIndex {
			// same edge case as above
			n.prevTail = prev.prevTail
		} else {
			panic("reached an unanticipated edge case! test better")
		}

		nextDirIndex := (n.pathIndex+1) * NameSize - startIndex
		if nextDirIndex < len(nodes) {
			n.nextDir = &nodes[nextDirIndex]
		}

		if i > 0 {
			n.prev = &nodes[i-1]
		} else {
			n.prev = prev
		}

		if i < len(nodes)-1 {
			n.next = []*PNode[T]{&nodes[i+1]}
		} else {
			// make it empty
			n.next = []*PNode[T]{}
		}
	}
	return nodes
}

// Returns the name of the parent directory (if it exists)
func (n PNode[T]) DirName() *Name {
	return n.dirPath.Base()
}

func (n PNode[T]) Dir() Path {
	return n.dirPath
}

// Returns the name of the current path (based on the first element)
//
// This is a "use at your own risk" sort of thing. It is possible, if you're
// not careful, that you're getting a non-unique name. I still have to think
// about what logic creates this, but the main thing is that generally this
// should only be called later in the tree.
func (n PNode[T]) Name() *Name {
	return n.name
}

func (n PNode[T]) IsRoot() bool {
	// The "root" of the subtree is the node that starts our
	// subtree. Since it's the start, it has no associated character and
	// thus is sort of special.
	return n.pathIndex == 0 && n.nameIndex == -1
}

func (n PNode[T]) IsLeaf() bool {
	return len(n.next) == 0
}

func (n PNode[T]) IsDir() bool {
	return len(n.next) > 1 || n.NextDir() != nil
}

func (n PNode[T]) IsFile() bool {
	return !n.IsDir()
}

func (n PNode[T]) IsRelativeRoot() bool {
	// return true if we're a node at the start of a new file/directory
	return n.nameIndex == 0
}

func (n *PNode[T]) PrevDir() *PNode[T] {
	return n.prevDir
}

func (n *PNode[T]) NextDir() *PNode[T] {
	return n.nextDir
}

func (pg *SubTree[T]) Leafs() []*PNode[T] {
	return pg.leafs
}

func (t SubTree[T]) Print() {
	t.root.print()
}

func (n PNode[T]) print() {
	if n.nameIndex >= 0 {
		if len(n.prev.next) > 1 {
			// this starts a split
			fmt.Println()
			for i := 0; i < (n.pathIndex * NameSize) + n.pathIndex + n.nameIndex - 1; i++ {
				fmt.Print(" ")
			}
		}
		if n.nameIndex == 0 {
			fmt.Print("/")
		}
		fmt.Printf("%x", n.char)
	}
	for _, next := range n.next {
		next.print()
	}
	if n.nameIndex == -1 {
		fmt.Println()
	}
}

func (n PNode[T]) ToPathObject() PathObject {
	return PathObject{}
}
