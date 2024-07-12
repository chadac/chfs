package chfs

import (
	"fmt"
	"sync"
)

type PNode struct {
	subtree *SubTree
	parent *PNode

	paths []*Path

	pathIndex int
	nameIndex int

	char byte
	terminal bool

	prev *PNode
	next []*PNode

	// points to the tail of the previous directory
	prevDir *PNode
	// points to the head of the next directory (if it exists)
	nextDir **PNode
}

// a PathGroup is a tree representing the combined representations of
// several paths
type SubTree struct {
	root *PNode
	paths []*Path
	leafs []*PNode
}

func newPNode(subtree *SubTree, parent *PNode, char byte, pIndex int, nIndex int, prevDir *PNode, nextDir **PNode) *PNode {
	n := PNode{}
	n.subtree = subtree
	n.parent = parent
	n.char = char
	n.paths = []*Path{}
	n.pathIndex = pIndex
	n.nameIndex = nIndex
	n.terminal = false
	n.prevDir = prevDir
	n.nextDir = nextDir
	return &n
}

// Generates the basic skeleton of a path.
func (n *PNode) populate(wg *sync.WaitGroup) {
	nextPathIndex := n.pathIndex
	nextNameIndex := n.nameIndex + 1
	if nextNameIndex == NameSize {
		// terminal node, move to the next
		nextPathIndex += 1
		nextNameIndex = 0
	}

	if nextPathIndex >= len(*n.paths[0]) {
		// terminal node
		n.next = nil
		for _, p1 := range n.paths {
			for i, p2 := range n.subtree.paths {
				if p1 == p2 {
					n.subtree.leafs[i] = n
				}
			}
		}

		if wg != nil {
			wg.Done()
		}
		return
	}

	c1 := (*n.paths[0])[nextPathIndex].Index(nextNameIndex)

	newPrevDir := n.prevDir
	if n.nameIndex == 0 {
		newPrevDir = n
	}

	newNextDir := n.nextDir
	if n.nameIndex == 0 || newNextDir == nil {
		*newNextDir = n
		newNextDir = new(*PNode)
	}

	var next []*PNode
	activeNext := 0

	if len(n.paths) == 1 {
		next = []*PNode{
			newPNode(n.subtree, n, c1, nextPathIndex, nextNameIndex, newPrevDir, newNextDir),
		}
		activeNext++
		next[0].paths = n.paths
	} else if len(n.paths) > 1 {
		split := false
		for _, p := range n.paths {
			if (*p)[nextPathIndex].Index(nextNameIndex) != c1 {
				split = true
				break
			}
		}
		if split {
			next = make([]*PNode, TreeSize)

			for _, p := range n.paths {
				c := (*p)[nextPathIndex].Index(nextNameIndex)

				// initialization of next node
				if next[c] == nil {
					next[c] = newPNode(n.subtree, n, c, nextPathIndex, nextNameIndex, newPrevDir, new(*PNode))
					activeNext++
				}

				next[c].paths = append(next[c].paths, p)
			}
		} else {
			next = []*PNode{
				newPNode(n.subtree, n, c1, nextPathIndex, nextNameIndex, newPrevDir, newNextDir),
			}
			activeNext++
			next[0].paths = n.paths
		}
	}

	n.next = next

	for _, n2 := range next {
		if n2 != nil {
			if wg != nil {
				go n2.populate(wg)
			} else {
				// paralellizing incurs too much overhead usually so... don't
				n2.populate(wg)
			}
		}
	}
}

func NewSubTree(paths []Path, parallel bool) *SubTree {
	newPaths := make([]*Path, len(paths))
	for i, _ := range paths {
		newPaths[i] = &paths[i]
	}

	g := SubTree{}
	g.paths = newPaths
	g.leafs = make([]*PNode, len(paths))

	root := newPNode(&g, nil, 0, 0, -1, nil, new(*PNode))
	root.paths = newPaths

	var wg *sync.WaitGroup = nil
	if parallel {
		wg = &sync.WaitGroup{}
		wg.Add(len(paths))
	}
	root.populate(wg)
	if wg != nil {
		wg.Wait()
	}

	return &g
}

// Returns the name of the parent directory (if it exists)
func (n *PNode) DirName() *string {
	if n.pathIndex <= 0 {
		return nil
	}
	return (*n.paths[0])[n.pathIndex-1].encoded
}

// Returns the name of the current path (if it exists)
//
// If this is non-unique this returns null.
func (n *PNode) Name() *string {
	if len(n.paths) > 1 {
		return (*n.paths[0])[n.pathIndex].encoded
	}
	return nil
}

func (n PNode) IsRoot() bool {
	// The "root" of the subtree is the node that starts our
	// subtree. Since it's the start, it has no associated character and
	// thus is sort of special.
	return n.pathIndex == 0 && n.nameIndex == -1
}

func (n PNode) IsLeaf() bool {
	return len(n.next) == 0
}

func (n PNode) IsDir() bool {
	return len(n.next) > 1 || n.NextDir() != nil
}

func (n PNode) IsFile() bool {
	return !n.IsDir()
}

func (n PNode) IsRelativeRoot() bool {
	// return true if we're a node at the start of a new file/directory
	return n.nameIndex == 0
}

func (n *PNode) PrevDir() *PNode {
	return n.prevDir
}

func (n *PNode) NextDir() *PNode {
	return *n.nextDir
}

// func (p *PNode) Leafs() []*PNode {
// 	if len(p.next) == 1 {
// 		return p.next[0].Leafs()
// 	} else {
// 		leafs := make([]*PNode, len(p.paths))
// 		i := 0
// 		for _, next := range p.next {
// 			if next != nil {
// 				newLeafs := next.Leafs()
// 				for j := 0; j < len(newLeafs); j++ {
// 					leafs[i] = newLeafs[j]
// 					i++
// 				}
// 			}
// 		}
// 		return leafs
// 	}
// }

func (pg *SubTree) Leafs() []*PNode {
	return pg.leafs
}

func (n PNode) print() {
	if n.nameIndex >= 0 {
		if len(n.parent.next) > 1 {
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
