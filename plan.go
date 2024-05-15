package vfs

type PlanRef struct {
	index byte
	node *PlanNode
}

type PlanNode struct {
	plan *Plan
	name *Name
	nameIndex int
	fullIndex int
	// if true, is the root of a new directory
	isDirRoot bool

	branch *Branch

	next []PlanRef
}

// A plan is a tree datastructure representing a write operation to
type Plan struct {
	rootNode *PlanNode
	vfs *VFS
	pathString string
	path *Path
	currentDir int
	currentNode int
}

func RunPlan(store Store, rootId *Checksum, paths *Path) (*Plan, error) {
	return nil, nil
}
