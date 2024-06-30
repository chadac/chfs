package chfs

// import (
// 	"fmt"
// 	"sync"
// )

type PlanNode struct {
	n PNode

	// now describe our write actions
	currentBranch *Branch
	branch *Branch
}

type Plan struct {
	pg PathGroup
	root PlanNode
}

// type PlanRootNode struct {
// 	branchKey *Checksum
// 	next PlanNode
// }

// type PlanRef struct {
// 	index byte
// 	node PlanNode
// 	skip *PlanNode
// }

// type PlanNode struct {
// 	// initialization variables
// 	name *Name
// 	dirIndex int
// 	nameIndex int
// 	// if this is a file, we write!
// 	file *FileEvent
// 	// this is the list of next items we'll visit
// 	next []PlanRef

// 	// generated from reset()
// 	branch *Branch

// 	// generated from plan()
// 	// links to the new parent branch
// 	newRef *BranchRef
// 	// links to the new branch generated
// 	newBranch *Branch
// }

// // A plan is a tree datastructure associated with a single write event
// //
// // We initially allocate the entire plan tree and then update it to reflect the
// // current state of the store with what we'd like to write.
// type Plan struct {
// 	paths []*Path
// 	rootNode PlanNode
// }


// func (n PlanRootNode) reset(store Store) error {
// 	branch, err := store.Branch(n.branchKey)
// 	if err != nil {
// 		return err
// 	}
// 	n.next.branch = branch
// 	return n.next.reset(store)
// }

// func (n PlanNode) plan() {
// 	// if we are a branch
// 	if n.file != nil {
// 		wg := sync.WaitGroup{}
// 		wg.Add(len(n.next))
// 		newBranch := *branch
// 		for _, ref := range n.next {
// 			go func() {
// 				defer wg.Done()
// 				if ref.skip != nil {
// 					ref.skip.plan()
// 					newBranch[ref.index] = ref.skip.newRef
// 				} else {
// 					ref.node.plan()
// 					newBranch[ref.index] = ref.node.newRef
// 				}
// 			}()
// 		}
// 		wg.Wait()
// 		n.newBranch = newBranch
// 		n.newRef := BranchRef{newBranch.checksum(),nil,nil,nil}
// 		if n.nameIndex == 0 {
// 			n.newRef.name = n.name
// 			n.newRef.isDir = true
// 		}
// 		// TODO: I need cleanup logic for when a delete causes a branch to be empty (i.e., stuff should be able to get simpler)
// 	} else { // otherwise we are a terminal (file) node
// 		n.newBranch = nil
// 		if n.file.file == nil {
// 			n.newRef = nil
// 		} else {
// 			n.newRef := BranchRef{n.file.file,n.name,false,n.file.executable}
// 		}
// 	}
// }

// func (ref PlanRef) reset(store Store, branch *Branch, skip *PlanNode) error {
// 	if branch != nil {
// 		bref := branch[ref.index]
// 		// we're creating a new subtree
// 		if bref == nil {
// 		}
// 	} else {
// 		// we're creating a new branch altogether
// 		branch = new(Branch{})
// 		ref.node.branch = branch
// 	}
// 	return ref.node.reset(store, skip)
// }

// func (n PlanNode) reset(store Store, passthru *PlanRef) error {
// 	wg := sync.WaitGroup{}
// 	wg.Add(len(n.next))

// 	nextPassthru := passthru
// 	if passthru != nil {
// 		if n.nameIndex == 0 {
// 			passthru.skip = &n
// 			nextPassthru = nil
// 		}
// 	}

// 	for _, ref := range n.next {
// 		go func() {
// 			defer wg.Done()
// 			branch, err := store.Branch(n.branch[ref.index].id)
// 			ref.node.branch = &branch
// 			ref.node.reset(store, nextPassthru)
// 		}()
// 	}
// 	wg.Wait()
// 	// if n.index == NameSize-1 {
// 	// 	branch, err := store.Branch(&n.passthru.id)
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}
// 	// 	for _, ref := range n.next {
// 	// 		go func() {
// 	// 			defer wg.Done()
// 	// 			if !n.passthru.name.equals(ref.node.getName()) {
// 	// 				ref.node.branch = nil
// 	// 				ref.node.parent = nil
// 	// 			} else {
// 	// 				ref.node.setBranch(branch)
// 	// 				ref.node.setPassthru(nil)
// 	// 			}
// 	// 			ref.node.reset(store)
// 	// 		}()
// 	// 	}
// 	// } else {
// 	// 	for _, ref := range n.next {
// 	// 		go func() {
// 	// 			defer wg.Done()
// 	// 			c := ref.index
// 	// 			var branch *Branch = nil
// 	// 			var passthru *Checksum = nil
// 	// 			if b.branch == nil {
// 	// 				// this means we're starting a brand new tree
// 	// 			} else if branchRef := b.branch[c]; branchRef == nil {
// 	// 				// error: unexpected terminal (shouldn't happen)
// 	// 				return fmt.Errorf("unexpected terminal input")
// 	// 			} else if branchRef.name != nil {
// 	// 				branch = nil
// 	// 				passthru = &branchRef
// 	// 			} else {
// 	// 				branch, err := store.Branch(branchRef.id)
// 	// 				if err != nil {
// 	// 					return err
// 	// 				}
// 	// 				passthru = nil
// 	// 			}
// 	// 			ref.node.setBranch(branch)
// 	// 			ref.node.setPassthru(passthru)
// 	// 			ref.node.reset(store)
// 	// 		}()
// 	// 	}
// 	// }
// 	wg.Wait()
// }

// func (n PlanFileNode) reset(store Store) error {
// 	return nil
// }

// func (n PlanRootNode) update() *BranchRef {
// 	ref := self.next.update()
// 	self.branchKey := ref.branch.checksum()
// 	return nil
// }

// func (n PlanDirNode) update() *BranchRef {
// 	if n.branch != nil {
// 		// update our branch
// 		for _, ref := range n.next {
// 			n.branch[ref.index] = ref.node.update()
// 		}

// 		newKey := n.branch.key()
// 		newRef := BranchRef{newKey,nil,nil}

// 		// if we're starting a new dir, indicate on the pointer
// 		if n.nameIndex == 0 {
// 			newRef.name = node.name
// 			newRef.isDir = true
// 		}

// 		// return the updated branch key
// 		return &newRef
// 	} else {
// 		if len(n.next) > 1 {
// 			// create new branch no matter what
// 		} else { // if len(n.next) == 1
// 			// pass it forward
// 			return n.ref[0].node.update()
// 		}
// 	}
// }

// func (n PlanFileNode) update() *BranchRef {
// 	ref := BranchRef{n.file, n.name, false}
// 	return &ref
// }

// func createPlanNode(plan *Plan, paths []*Path, files []*Checksum, level int, char int) PlanNode {
// 	if level == 0 && char == -1 {
// 		next := CreatePlanNode(paths, files, 0, 0)
// 		node := PlanRootNode{nil, next}
// 		return &node
// 	} else {
// 		// todo: if everything is identical we shouldn't realloc
// 		pmap := make(map[int][]*Path)
// 		fmap := make(map[int][]*Checksum)
// 		for p := 0; i < len(paths); p++ {
// 			c := paths[p][level][char]
// 			pmap[c] = append(pmap[c], paths[p])
// 			fmap[c] = append(fmap[c], files[p])
// 		}
// 		next := make([]PlanRef)
// 		for c, subpaths := range pmap {
// 			nextLevel := level
// 			nextChar := char+1
// 			if nextChar >= NameSize {
// 				nextLevel += 1
// 				nextChar = 0
// 			}
// 			nextNode := createPlanNode(plan, subpaths, fmap[c], level, char+1)
// 			ref := PlanRef{c, nextNode}
// 			next = append(next, &ref)
// 		}
// 		node := PlanNode{plan, level, index, nil, next}
// 	}
// }

// func CreateEmptyPlan(paths []*Path, files []*Checksum) *Plan {
// 	plan := Plan{}
// 	plan.store = store
// 	plan.paths = [1]*Path{path}
// 	rootNode := PlanRootNode{&plan, root}
// 	var currNode PlanNode
// 	for i, name := range path {
// 		for j := 0; j < NameSize; j++ {
// 			c := name.encoded.index(j)
// 		}
// 	}
// 	return &plan, nil
// }
