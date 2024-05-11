package vfs

type planRef struct {
	index byte
	node planNode
}

type planNode struct {
	plan *plan
	name *Name
	nameIndex int
	fullIndex int
	// if true, is the root of a new directory
	isDirRoot bool

	branch *Branch

	next []planRef
}

// A plan is a tree datastructure representing a write operation to
type plan struct {
	rootNode *planNode
	vfs *VFS
	pathString string
	path *Path
	currentDir int
	currentNode int
}

func (n routeNode) key() string {
	return n.route.pathString[:n.index]
}

func (n routeNode) next() *branchNode {
	return n.branch[n.nextChar]
}

func (d routeDir) nextName() *Name {
	return nil
}

func Plan(store Store, rootId *Checksum, paths *Path)
