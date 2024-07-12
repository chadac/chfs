package chfs

import (
	"sync"
)

type PlanNode struct {
	plan *Plan
	n *PNode

	// read attributes
	treeId chan *Checksum
	currTree *Tree
	currFile *File

	// stuff we write
	newTree *Tree
	newFile *File

	// mimic structural stuff
	next []*PlanNode
	prev *PlanNode
	prevDir *PlanNode
	nextDir *PlanNode
}

type Plan struct {
	pg *SubTree
	root *PlanNode
	leafs []*PlanNode
}

func NewPlan(st *SubTree) *Plan {
	plan := Plan{st,nil,make([]*PlanNode, len(st.leafs))}
	root := newPlanNode(&plan, st.root)
	var nextDir *PlanNode = nil
	if root.n.NextDir() != nil {
		nextDir = newPlanNode(&plan, root.n.NextDir())
	}
	root.init(nextDir)
	return &plan
}

func newPlanNode(plan *Plan, n *PNode) *PlanNode {
	pn := new(PlanNode)
	pn.plan = plan
	pn.n = n
	pn.treeId = make(chan *Checksum)
	pn.next = make([]*PlanNode, len(n.next))
	return pn
}

func (pn *PlanNode) init(nextDir *PlanNode) {
	pn.nextDir = nextDir
	n := pn.n
	for i, next := range n.next {
		var pnext *PlanNode
		if nextDir != nil && nextDir.n == next {
			pnext = nextDir
		} else {
			pnext = newPlanNode(pn.plan, next)
		}
		pn.next[i] = pnext
		pnext.prev = pn

		if pn.prevDir != nil && pn.prevDir.n == next.prevDir {
			pnext.prevDir = pn.prevDir
		} else if next.prevDir != nil {
			// this means we're swapping directories... reassign properly
			pnext.prevDir = pn
		}

		if nextDir != nil && nextDir.n == next.NextDir() {
			pnext.init(nextDir)
		} else if next.NextDir() != nil {
			newNextDir := newPlanNode(pn.plan, next.NextDir())
			pnext.init(newNextDir)
		} else {
			pnext.init(nil)
		}
	}
}

func (p *Plan) Read(rootId Checksum, store Store[Checksum, Tree]) {
	wg := sync.WaitGroup{}
	go p.root.read(store, &wg)
	p.root.treeId <- &rootId
	wg.Wait()
}

func (pn *PlanNode) read(store Store[Checksum, Tree], wg *sync.WaitGroup) {
	for _, nn := range pn.next {
		go nn.read(store, wg)
	}

	go func() {
		treeId := <-pn.treeId
		if treeId != nil {
			// TODO: figure out read errors
			tree, _ := store.Get(*treeId)
			pn.currTree = &tree
			for _, nn := range pn.next {
				nn.readBranch(tree.b[nn.n.char])
			}
		}
		if pn.n.IsLeaf() {
			wg.Add(1)
		}
	}()
}

func (pn *PlanNode) readBranch(branch *Branch) {
	if branch == nil {
		pn.treeId <- nil
	} else if branch.IsTerminal() {
		if pn.nextDir == nil {
			// this means that we're going to need to rewrite...
			pn.treeId <- nil
		} else {
			nextDir := pn.nextDir
			if *pn.n.Name() == branch.obj.name {
				nextDir.treeId <- &branch.id
				if nextDir != pn {
					pn.treeId <- nil
				}
			} else {
				pn.treeId <- nil
			}
		}
	} else if branch.IsFile() {
		pn.currFile = &File{branch.id}
		pn.treeId <- nil
	} else {
		// otherwise we just send the object down
		pn.treeId <- &branch.id
	}
}

func (p *Plan) Update(files []File) {
	for i, leaf := range p.leafs {
		leaf.newFile = &files[i]
	}

	// generate the new tree
	p.root.updateTree()
}

func (pn PlanNode) findNewFile() *File {
	if len(pn.next) > 1 {
		return nil
	}
	return pn.next[0].findNewFile()
}

func (pn *PlanNode) updateTree() {
	// Using the read value, update the present tree to whatever the new
	// thing we want is.

	// TODO: This may need some options for update strategies.

	var newTree *Tree = nil
	if pn.currTree != nil {
		// simple case: if we're running on an existing tree, then we simply run on existing files
		newTree = CopyTree(pn.currTree)
	} else {
		// otherwise we're generating a new tree
		newTree = EmptyTree()
	}

	// now we populate our new tree with all the new branches
	for _, next := range pn.next {
		newTree.b[next.n.char] = next.createBranch()
	}

	numBranches := newTree.BranchCount()

	// consolidation: if our tree is empty, we'll delete it
	if numBranches == 0 {
		newTree = nil
	}

	// consolidation: if our tree is non-relroot and has only one branch, consolidate
	if !pn.n.IsRelativeRoot() && numBranches <= 1 {
		newTree = nil
	}

	pn.newTree = newTree
}

func (pn *PlanNode) createBranch() *Branch {
	// TODO: We may need to conditionalize this for partial tree updates
	pn.updateTree()

	// if we can consolidate, then just jump to the next tree
	if pn.newTree == nil && pn.n.IsDir() {
		// case 1: we don't need to create a new tree but there is a new
		// directory
		return pn.nextDir.createBranch()
	} else if pn.newTree == nil { // && pn.n.IsFile()
		// case 2: if the tree has been "consolidated" that means that we don't need
		// to split anymore. get the file at the end of the rainbow and return it
		pn.newFile = pn.findNewFile()
		if pn.newFile == nil {
			// PANIC!
		}
		obj := NewFile(*pn.n.Name())
		return &Branch{pn.newFile.Key(), &obj}
	} else if pn.n.nameIndex == 0 {
		// case three: (complement to case one) we're at the root of a new directory
		obj := NewDir(*pn.n.Name())
		return &Branch{pn.newTree.Key(),&obj}
	} else {
		// case four: regular case. we're pointing to the tree that this node represents
		return &Branch{pn.newTree.Key(),nil}
	}
}

func (p *Plan) Write(store Store[Checksum, Tree]) *Tree {
	wg := sync.WaitGroup{}
	p.root.write(&wg, store)
	wg.Wait()

	// return the new root tree
	return p.root.newTree
}

func (pn *PlanNode) write(wg *sync.WaitGroup, store Store[Checksum, Tree]) {
	if pn.currTree == nil {
		if pn.nextDir != nil {
			pn.nextDir.write(wg, store)
		}
		return
	}
	wg.Add(1)
	for _, next := range pn.next {
		next.write(wg, store)
	}
	go func() {
		defer wg.Done()
		if pn.currTree != nil {
			store.Put(*pn.currTree)
		}
	}()
}
