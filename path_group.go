package chfs

import (
	"fmt"
)

type PNode struct {
	parent *PNode

	paths []Path

	pathIndex int
	nameIndex int

	char byte
	terminal bool

	next []*PNode

	read *RNode
	write *WNode
}

type RNode struct {
	branch *Branch
	tree *Tree

	// transitive parent
	parent *PNode
	// trasitive child
	child *PNode
}

type WNode struct {
	parent *Branch
	tree *Tree
}

// a PathGroup is a tree representing the combined representations of
// several paths
type PathGroup struct {
	root *PNode
	paths *[]Path
	leafs []*PNode
}

func createPNode(parent *PNode, char byte, pIndex int, nIndex int) *PNode {
	n := PNode{}
	n.parent = parent
	n.char = char
	n.paths = []Path{}
	// n.next = []*PNode{}
	n.pathIndex = pIndex
	n.nameIndex = nIndex
	n.terminal = false
	return &n
}

// Generates the basic skeleton of a path.
func (n *PNode) populate() {
	// fmt.Printf("start %d %d\n", n.pathIndex, n.nameIndex)

	nextPathIndex := n.pathIndex
	nextNameIndex := n.nameIndex + 1
	if nextNameIndex == NameSize {
		// terminal node, move to the next
		nextPathIndex += 1
		nextNameIndex = 0
	}

	if nextPathIndex >= len(n.paths) {
		// terminal node
		n.next = nil
		return
	}

	c1 := n.paths[0][nextPathIndex].Index(nextNameIndex)

	var next []*PNode
	activeNext := 0

	if len(n.paths) == 1 {
		next = []*PNode{
			createPNode(n, c1, nextPathIndex, nextNameIndex),
		}
		activeNext++
		next[0].paths = n.paths
	} else if len(n.paths) > 1 {
		split := false
		for _, p := range n.paths {
			if p[nextPathIndex].Index(nextNameIndex) != c1 {
				split = true
				break
			}
		}
		if split {
			next = make([]*PNode, TreeSize)

			for _, p := range n.paths {
				c := p[nextPathIndex].Index(nextNameIndex)

				// initialization of next node
				if next[c] == nil {
					next[c] = createPNode(n, c, nextPathIndex, nextNameIndex)
					activeNext++
				}

				next[c].paths = append(next[c].paths, p)
			}
		} else {
			next = []*PNode{
				createPNode(n, c1, nextPathIndex, nextNameIndex),
			}
			activeNext++
			next[0].paths = n.paths
		}
	}

	n.next = next

	// if activeNext > 8 {
	// 	fmt.Println("hmm")
	// 	var wg sync.WaitGroup
	// 	for _, n2 := range next {
	// 		if n2 != nil {
	// 			wg.Add(1)
	// 			// wg.Done()
	// 			go func(node *PNode) {
	// 				defer wg.Done()
	// 				node.populate()
	// 			}(n2)
	// 		}
	// 	}
	// 	wg.Wait()
	// } else {

	// paralellizing incurs too much overhead tbh
	for _, n2 := range next {
		if n2 != nil {
			n2.populate()
		}
	}

	// fmt.Printf("done %d %d\n", n.pathIndex, n.nameIndex)
}

// This
func (n *PNode) prevName() *string {
	if n.pathIndex <= 0 {
		return nil
	}
	return n.paths[0][n.pathIndex-1].encoded
}

// Loads a tree with the existing plan
// structure:
// [ /    ] RNode { tree = /, branch = nil }
// [ /a   ] RNode { tree = /a, branch = /->a }
// [ /ab  ] RNode { tree = /ab, branch = /a->b }
// [ /abc ] RNode { tree = nil, branch = /ab->c }
func (n *PNode) load(store Store[Checksum, Tree], carry *PNode) error {
	var read *RNode = nil
	var treeId *Checksum = nil
	if carry != nil && n.nameIndex == 0 {
		if carry.read.branch == nil {
			return fmt.Errorf(`expected carry to contain a branch`)
		}
		obj := carry.read.branch.obj
		// TODO: directory check?
		if obj == nil || obj.name != *n.prevName() {
			// this means that the branch that the parent tree is pointing
			// from broke everything
			return nil
		}
		read = &RNode{}
		carry.read.child = n
		read.parent = carry
		treeId = &carry.read.branch.id
		carry = nil
	} else if carry == nil {
		read = &RNode{}
		read.branch = n.parent.read.tree.b[n.char]
		if read.branch == nil {
			// this means we're creating a new directory, we can cancel here
			return nil
		} else if read.branch.obj != nil {
			carry = read.parent
		} else {
			treeId = &read.branch.id
		}
	}

	if treeId != nil {
		tree, err := store.Get(*treeId)
		if err != nil {
			return err
		}
		read.tree = &tree
	}

	// populate children
	// TODO: parallelization should be a flag in case something is IO-bound
	for _, next := range n.next {
		next.load(store, carry)
	}

	return nil
}

func CreatePathGroup(paths []Path) *PathGroup {
	root := createPNode(nil, 0, 0, -1)
	root.paths = paths
	root.populate()
	g := PathGroup{root,&paths,nil}
	return &g
}

func (p *PNode) Leafs() []*PNode {
	if len(p.next) == 1 {
		return p.next[0].Leafs()
	} else {
		leafs := make([]*PNode, len(p.paths))
		i := 0
		for _, next := range p.next {
			if next != nil {
				newLeafs := next.Leafs()
				for j := 0; j < len(newLeafs); j++ {
					leafs[i] = newLeafs[j]
					i++
				}
			}
		}
		return leafs
	}
}

func (pg *PathGroup) Leafs() []*PNode {
	if pg.leafs == nil {
		pg.leafs = pg.root.Leafs()
	}
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
	// if len(n.next) <= 0 {
	// 	fmt.Println()
	// }
	for _, next := range n.next {
		next.print()
	}
	if n.nameIndex == -1 {
		fmt.Println()
	}
}
