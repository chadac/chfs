package vfs

import (
	"fmt"
	"sync"
)

type PlanRootNode struct {
	branchKey *Checksum
	next PlanNode
}

type PlanRef struct {
	index byte
	node PlanNode
}

type PlanNode struct {
	plan *Plan
	name *Name
	level int
	index int

	// we edit this one
	branch *Branch
	// this is used to propagate changes
	passthru *BranchRef

	isLeaf bool
	file *Checksum

	next []PlanRef
}

// A plan is a tree datastructure representing a write operation to
type Plan struct {
	paths []*Path
	rootNode PlanNode
}

func (n PlanRootNode) reset(store Store) error {
	branch, err := store.Branch(n.branchKey)
	if err != nil {
		return err
	}
	n.next.branch = branch
	return n.next.reset(store)
}

func (n PlanNode) reset(store Store) error {
	wg := sync.WaitGroup{}
	wg.Add(len(n.next))
	if n.index == NameSize-1 {
		branch, err := store.Branch(&n.passthru.id)
		if err != nil {
			return err
		}
		for _, ref := range n.next {
			if !n.passthru.name.equals(ref.node.getName()) {
				ref.node.setBranch(nil)
				ref.node.setPassthru(nil)
			} else {
				ref.node.setBranch(branch)
				ref.node.setPassthru(nil)
			}
			go func() {
				defer wg.Done()
				ref.node.reset(store)
			}()
		}
	} else {
		for i, ref := range n.next {
			c := ref.index
			var branch *Branch = nil
			var passthru *Checksum = nil
			if b.branch == nil {
				// this means we're starting a brand new tree
			} else if branchRef := b.branch[c]; branchRef == nil {
				// error: unexpected terminal (shouldn't happen)
				return fmt.Errorf("unexpected terminal input")
			} else if branchRef.name != nil {
				branch = nil
				passthru = &branchRef
			} else {
				branch, err := store.Branch(branchRef.id)
				if err != nil {
					return err
				}
				passthru = nil
			}
			ref.node.setBranch(branch)
			ref.node.setPassthru(passthru)
			go func() {
				ref.node.reset(store) 
				defer wg.Done()
			}()
		}
	}
	wg.Wait()
}

func (n PlanFileNode) reset(store Store) error {
	return nil
}

func (n PlanRootNode) update() *BranchRef {
	ref := self.next.update()
	self.branchKey := ref.branch.checksum()
	return nil
}

func (n PlanDirNode) update() *BranchRef {
	if n.branch != nil {
		// update our branch
		for _, ref := range n.next {
			n.branch[ref.index] = ref.node.update()
		}

		newKey := n.branch.key()
		newRef := BranchRef{newKey,nil,nil}

		// if we're starting a new dir, indicate on the pointer
		if n.nameIndex == 0 {
			newRef.name = node.name
			newRef.isDir = true
		}

		// return the updated branch key
		return &newRef
	} else {
		if len(n.next) > 1 {
			// create new branch no matter what
		} else { // if len(n.next) == 1
			// pass it forward
			return n.ref[0].node.update()
		}
	}
}

func (n PlanFileNode) update() *BranchRef {
	ref := BranchRef{n.file, n.name, false}
	return &ref
}

func createPlanNode(plan *Plan, paths []*Path, files []*Checksum, level int, char int) PlanNode {
	if level == 0 && char == -1 {
		next := CreatePlanNode(paths, files, 0, 0)
		node := PlanRootNode{nil, next}
		return &node
	} else {
		// todo: if everything is identical we shouldn't realloc
		pmap := make(map[int][]*Path)
		fmap := make(map[int][]*Checksum)
		for p := 0; i < len(paths); p++ {
			c := paths[p][level][char]
			pmap[c] = append(pmap[c], paths[p])
			fmap[c] = append(fmap[c], files[p])
		}
		next := make([]PlanRef)
		for c, subpaths := range pmap {
			nextLevel := level
			nextChar := char+1
			if nextChar >= NameSize {
				nextLevel += 1
				nextChar = 0
			}
			nextNode := createPlanNode(plan, subpaths, fmap[c], level, char+1)
			ref := PlanRef{c, nextNode}
			next = append(next, &ref)
		}
		node := PlanNode{plan, level, index, nil, next}
	}
}

func CreateEmptyPlan(paths []*Path, files []*Checksum) *Plan {
	plan := Plan{}
	plan.store = store
	plan.paths = [1]*Path{path}
	rootNode := PlanRootNode{&plan, root}
	var currNode PlanNode
	for i, name := range path {
		for j := 0; j < NameSize; j++ {
			c := name.encoded.index(j)
		}
	}
	return &plan, nil
}
