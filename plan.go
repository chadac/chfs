package chfs

import (
	_ "fmt"
	"sync"
)

type readData interface {
	TreeId() chan *Checksum
	CurrTree() *Tree
	CurrFile() *Branch

	sendTreeId(*Checksum)
	setCurrTree(*Tree)
	setCurrFile(*Branch)
}

type writeData interface {
	readData

	NewTree() *Tree
	NewFile() *Branch

	setNewTree(*Tree)
	setNewFile(*Branch)
}

type Reader struct {
	treeId chan *Checksum
	currTree *Tree
	currFile *Branch
}

func NewReader() *Reader {
	return &Reader{
		make(chan *Checksum),
		nil,
		nil,
	}
}

func (r Reader) TreeId() chan *Checksum { return r.treeId }
func (r Reader) CurrTree() *Tree { return r.currTree }
func (r Reader) CurrFile() *Branch { return r.currFile }
func (r *Reader) sendTreeId(treeId *Checksum) { r.treeId <- treeId }
func (r *Reader) setCurrTree(tree *Tree) { r.currTree = tree }
func (r *Reader) setCurrFile(obj *Branch) { r.currFile = obj }

type Writer struct {
	treeId chan *Checksum
	currTree *Tree
	currFile *Branch
	newTree *Tree
	newFile *Branch
}

func NewWriter() *Writer {
	return &Writer{
		make(chan *Checksum),
		nil,
		nil,
		nil,
		nil,
	}
}

func (w Writer) TreeId() chan *Checksum { return w.treeId }
func (w Writer) CurrTree() *Tree { return w.currTree }
func (w Writer) CurrFile() *Branch { return w.currFile }
func (w Writer) NewTree() *Tree { return w.newTree }
func (w Writer) NewFile() *Branch { return w.newFile }
func (w *Writer) sendTreeId(treeId *Checksum) { w.treeId <- treeId }
func (w *Writer) setCurrTree(tree *Tree) { w.currTree = tree }
func (w *Writer) setNewTree(tree *Tree) { w.newTree = tree }
func (w *Writer) setCurrFile(obj *Branch) { w.currFile = obj }
func (w *Writer) setNewFile(obj *Branch) { w.newFile = obj }

func ReadTree[T readData](st *SubTree[T], store TreeStore, rootId Checksum) {
	wg := sync.WaitGroup{}
	go readNode(st.root, store, &wg)
	// fmt.Printf("%p\n\n", &st.root)
	// fmt.Printf("%p\n\n", &st.root.extra)
	st.root.extra.sendTreeId(&rootId)
	wg.Wait()
}

func readNode[T readData](n *PNode[T], store TreeStore, wg *sync.WaitGroup) {
	for _, next := range n.next {
		go readNode(next, store, wg)
	}

	go func() {
		treeId := <-n.extra.TreeId()
		if treeId != nil {
			// TODO: figure out read errors
			tree, _ := store.Get(*treeId)
			n.extra.setCurrTree(tree)

			for _, next := range n.next {
				readNodeBranch(next, tree.b[next.char])
			}
		}

		if n.IsLeaf() {
			wg.Add(1)
		}
	}()
}

func readNodeBranch[T readData](n *PNode[T], branch *Branch) {
	if branch == nil {
		n.extra.sendTreeId(nil)
	} else if branch.IsTerminal() {
		if n.nextDir == nil {
			// this means that we're going to need to rewrite...
			n.extra.sendTreeId(nil)
		} else {
			nextDir := n.NextDir()
			if n.Name().Encoded() == branch.obj.name.Encoded() {
				nextDir.extra.sendTreeId(&branch.id)
				if nextDir != n {
					n.extra.sendTreeId(nil)
				}
			} else {
				n.extra.sendTreeId(nil)
			}
		}
	} else if branch.IsFile() {
		n.extra.setCurrFile(branch)
		n.extra.sendTreeId(nil)
	} else {
		// otherwise we just send the object down
		n.extra.sendTreeId(&branch.id)
	}
}

func Plan[T writeData](t *SubTree[T], files []*Branch) {
	for i, leaf := range t.leafs {
		leaf.extra.setNewFile(files[i])
	}

	// generate the new tree
	updateTree(t.root)
}

func findNewFile[T writeData](n *PNode[T]) *Branch {
	if n.extra.NewFile() != nil {
		return n.extra.NewFile()
	}
	if len(n.next) != 1 {
		return nil
	}
	return findNewFile(n.next[0])
}

func updateTree[T writeData](n *PNode[T]) {
	// Using the read value, update the present tree to whatever the new
	// thing we want is.

	// TODO: This may need some options for update strategies.

	var newTree *Tree = nil
	if n.extra.CurrTree() != nil {
		// simple case: if we're running on an existing tree, then we simply run on existing files
		newTree = CopyTree(n.extra.CurrTree())
	} else {
		// otherwise we're generating a new tree
		newTree = EmptyTree()
	}

	// now let's count the number of branches we're going to traverse
	totalBranches := newTree.BranchCount()
	for _, next := range n.next {
		if newTree.b[next.char] == nil {
			totalBranches++
		}
	}

	// consolidation: if our new tree is guaranteed to be simple, then don't make it
	if !n.IsRoot() && !n.IsRelativeRoot() && totalBranches <= 1 {
		// fmt.Printf("%s: not enough branches\n", n)
		newTree = nil
	} else {
		// we need to split; therefore
		for _, next := range n.next {
			newTree.b[next.char] = createBranchFor(next)
		}

		numBranches := newTree.BranchCount()

		if numBranches == 0 {
			// consolidation: if our tree is empty, we'll delete it
			// fmt.Printf("%s: not enough branches after creation\n", n)
			newTree = nil
		} else if !n.IsRoot() && !n.IsRelativeRoot() && numBranches <= 1 {
			// consolidation: if our tree is non-relroot and has only one branch, consolidate
			// fmt.Printf("%s: tree too small, consolidating\n", n)
			newTree = nil
		}
	}

	// fmt.Printf("%s: new tree %p\n", n, newTree)
	n.extra.setNewTree(newTree)
}

func createBranchFor[T writeData](n *PNode[T]) *Branch {
	// TODO: We may need to conditionalize this for partial tree updates
	updateTree(n)

	// if we can consolidate, then just jump to the next tree
	if n.extra.NewTree() == nil && n.IsDir() {
		// case 1: we don't need to create a new tree but there is a new
		// directory
		return createBranchFor(n.NextDir())
	} else if n.extra.NewTree() == nil { // && pn.n.IsFile()
		// case 2: if the tree has been "consolidated" that means that we don't need
		// to split anymore. get the file at the end of the rainbow and return it
		newFile := findNewFile(n)
		n.extra.setNewFile(newFile)
		if n.extra.NewFile() == nil {
			// PANIC!
		}
		// TODO: FIX
		return newFile
	} else if n.nameIndex == 0 {
		// case three: (complement to case one) we're at the root of a new directory
		obj := NewDir(*n.Name())
		return &Branch{n.extra.NewTree().Key(),&obj}
	} else {
		// case four: regular case. we're pointing to the tree that this node represents
		return &Branch{n.extra.NewTree().Key(),nil}
	}
}

func WriteTree[T writeData](t *SubTree[T], store Store[Checksum, Tree]) *Tree {
	wg := sync.WaitGroup{}
	writeNode(t.root, &wg, store)
	wg.Wait()

	// return the new root tree
	return t.root.extra.NewTree()
}

func writeNode[T writeData](n *PNode[T], wg *sync.WaitGroup, store Store[Checksum, Tree]) {
	if n.extra.CurrTree() == nil {
		if n.NextDir() != nil {
			writeNode(n.NextDir(), wg, store)
		}
		return
	}
	wg.Add(1)
	for _, next := range n.next {
		writeNode(next, wg, store)
	}

	// async write the tree to the store
	go func() {
		defer wg.Done()
		if n.extra.NewTree() != nil {
			store.Put(*n.extra.NewTree())
		}
	}()
}
