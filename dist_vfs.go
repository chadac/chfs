package vfs

// type DistVFS struct {
// 	store Store
// 	cache WriteCache
// }

// func (vfs *VFS) Reset() (*Checksum, error) {
// 	emptyBranch := Branch{}
// 	key, err := vfs.store.Put(emptyBranch)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return vfs.setRoot(key)
// }

// func (vfs *VFS) Root() (*Checksum, error) {
// 	rootObj, err := vfs.store.Get(&rootKey)
// 	if err != nil {
// 		return nil, err
// 	}
// 	rootFile, ok := rootObj.(*File)
// 	if !ok {
// 		return nil, fmt.Errorf("root object not of file type")
// 	}
// 	return (*Checksum)(rootFile), nil
// }

// func (vfs *VFS) setRoot(branchId *Checksum) (*Checksum, error) {
// 	err := vfs.store.Set(&rootKey, (*File)(branchId))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &rootKey, nil
// }

// func (vfs *VFS) Branch(id *Checksum) (*Branch, error) {
// 	obj, err := vfs.store.Get(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	branch, ok := obj.(Branch)
// 	if !ok {
// 		return nil, fmt.Errorf("object at location '%x' is not a branch", id)
// 	}
// 	return &branch, nil
// }

// func (vfs *VFS) File(id *Checksum) (*File, error) {
// 	obj, err := vfs.store.Get(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	file, ok := obj.(*File)
// 	if !ok {
// 		return nil, fmt.Errorf("object at location '%x' is not a file", id)
// 	}
// 	return file, nil
// }

// func (vfs *VFS) Get(path *Path) (*File, error) {
// 	curr, err := vfs.Root()
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, name := range *path {
// 		for j := 0; j < 32; j++ {
// 			b, err := vfs.Branch(curr)
// 			if err != nil {
// 				return nil, err
// 			}
// 			p := b.next(name.index(j))
// 			if p == nil {
// 				// we've reached a terminal node
// 				return nil, fmt.Errorf("directory does not exist")
// 			}
// 			if p.name != nil {
// 				if name.equals(p.name) {
// 					break
// 				} else {
// 					return nil, fmt.Errorf("directory does not exist")
// 				}
// 			}
// 			curr = &p.id
// 		}
// 	}

// 	return vfs.File(curr)
// }

// type routeDir struct {
// 	route *route
// 	index int
// 	n []*routeNode
// 	name *Name
// 	start int
// 	needsSplit bool
// }

// type routeNode struct {
// 	route *route
// 	dir *routeDir
// 	branch *Branch
// 	nextChar byte
// 	index int
// 	dirIndex int
// }

// type route struct {
// 	d []*routeDir
// 	vfs *VFS
// 	pathString string
// 	path *Path
// 	currentDir int
// 	currentNode int
// }

// func (n routeNode) key() string {
// 	return n.route.pathString[:n.index]
// }

// func (n routeNode) next() *branchNode {
// 	return n.branch[n.nextChar]
// }

// func (d routeDir) nextName() *Name {
// 	return nil
// }

// func (d routeDir) isLast() bool {
// 	return d.index == len(d.route.d) - 1
// }

// func (d *routeDir) populate(from int) error {
// 	var err error
// 	vfs := d.route.vfs
// 	for i := from; i < nameSize-1; i++ {
// 		fmt.Printf("populate idx: %d\n", i)
// 		next := d.n[i].next()
// 		if next == nil {
// 			d.start = i
// 			d.needsSplit = true
// 			return nil
// 		} else if next.name == nil {
// 			d.n[i+1].branch, err = vfs.Branch(&next.id)
// 			if err != nil {
// 				return err
// 			}
// 		} else if next.name.equals(d.nextName()) {
// 			// this means we're overwriting whatever existed here originally
// 			d.start = i
// 			return nil
// 		} else {
// 			// we're going to need to create a directory
// 			d.start = i
// 			d.needsSplit = true
// 			return nil
// 		}
// 	}
// 	return fmt.Errorf("directory does not terminate; improperly formatted?")
// }

// func (d routeDir) nextDir() *Checksum {
// 	last := d.n[d.start]
// 	bn := last.branch[last.nextChar]
// 	return &bn.id
// }

// func (r *route) next() *routeDir {
// 	r.currentDir -= 1
// 	if r.currentDir < 0 {
// 		return nil
// 	}
// 	return r.d[r.currentDir]
// }

// // used for building plans to run an update
// func (vfs VFS) route(path *Path) (*route, error) {
// 	r := route{}
// 	r.vfs = &vfs
// 	r.pathString = path.encoded()
// 	r.path = path
// 	// last name is the filename... second to last is the last dirname
// 	r.currentDir = len(*path)-2
// 	r.d = make([]*routeDir, len(*path))

// 	prev, err := vfs.Root()
// 	if err != nil {
// 		return nil, err
// 	}

// 	fmt.Printf("%x\n", *prev)
// 	skip := false
// 	for i := 0; i < len(*path); i++ {
// 		fmt.Printf("%d\n", i)
// 		d := routeDir{}
// 		d.route = &r
// 		d.index = i
// 		d.name = &(*path)[i]
// 		d.n = make([]*routeNode, nameSize)
// 		d.needsSplit = false
// 		for j := 0; j < nameSize; j++ {
// 			n := routeNode{}
// 			n.route = &r
// 			n.dir = &d
// 			n.branch = nil
// 			n.nextChar = (*path)[i].index(j)
// 			n.index = i * nameSize + j
// 			n.dirIndex = j
// 			d.n[j] = &n
// 		}
// 		if !skip {
// 			d.n[0].branch, err = vfs.Branch(prev)
// 			if err != nil {
// 				return nil, err
// 			}
// 			err = d.populate(0)
// 			if err != nil {
// 				return nil, err
// 			}
// 		} else {
// 			d.n[0].branch = &Branch{}
// 		}
// 		if d.needsSplit {
// 			skip = true
// 		} else {
// 			prev = d.nextDir()
// 		}
// 		r.d[i] = &d
// 	}
// 	return &r, nil
// }

// func (vfs *VFS) Gets(paths []*Path) ([]*File, error) {
// 	// TODO
// 	return nil, nil
// }

// func (vfs *VFS) Set(path *Path, file *File) (*Checksum, error) {
// 	// then generate the planned route for updating
// 	route, err := vfs.route(path)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// start by claiming all subpaths
// 	id := vfs.cache.Claim(route.pathString)

// 	curr, err := vfs.store.Put(file)
// 	if err != nil {
// 		return nil, err
// 	}

// 	nextEntry := &branchNode{*curr,path.fileName(),false}

// 	// iterate from the tail of the tree to the start
// 	for dir := route.next(); dir != nil; {
// 		fmt.Printf("directory: %s", dir.name.repr)
// 		for j := nameSize-1; j >= 0; j-- {
// 			n := dir.n[j]
// 			if newB := vfs.cache.Acquire(n.key(), id); newB != nil {
// 				n.branch = newB
// 			}
// 			if n.branch != nil {
// 				n.branch[n.nextChar] = nextEntry
// 			}
// 			curr, err = vfs.store.Put(n.branch)
// 			nextEntry = &branchNode{*curr,nil,false}
// 			// oh no look at that!!! if an error happens we don't clean up
// 			// what's been written so far!!!!
// 			if err != nil {
// 				// TODO: release all our other locks... we failed :(
// 				return nil, err
// 			}
// 			// at the end of this whole thing... we gotta point to the new root
// 			if j == 0 && dir.index == 0 {
// 				_, err = vfs.setRoot(curr)
// 				if err != nil {
// 					// todo release the rest
// 					return nil, err
// 				}
// 			}
// 			vfs.cache.Release(n.key(), n.branch)
// 		}
// 		nextEntry.name = dir.name
// 		nextEntry.isDir = true
// 	}

// 	return curr, nil
// }

// func (vfs *VFS) Sets(paths []*Path, files []*File) ([]*Checksum, error) {
// 	// BIG TODO
// 	// this is a whole lot harder...
// 	return nil, nil
// }
