package main

type path []pathRef

type pathRef struct {
	root *id
	subpath []*id
	next *id
}

type branch [baseSize]branchRef

type interface branchRef {
	name() *checksum
	next() *id
}

type bBranch {
	next *id
}

type bLink {
	id *checksum
	next *checksum
}

func (b *branch) get(key byte) (*branchRef, error) {
	if key > baseSize {
		return nil, fmt.Errorf("index '%+v' out of range", key)
	}
	return b[key], nil
}

func (b *branch) update(key byte, newRef *branchRef) (*branch, error) {
	newBranch := make(branch)
	copy(newBranch, b)
	if key > baseSize {
		return nil, fmt.Errorf("index '%+v' out of range", key)
	}
	newBranch[key] = newRef
	return &newBranch
}
