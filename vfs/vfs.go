package vfs

import (
	"fmt"
)

// The root is stored at a single predictable location
// It is the one mutable part of the store. This is used to track the root of
// the trie for consistency.
var rootKey = Checksum{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

type VFS struct {
	store Store
	cache WriteCache
}

func (vfs *VFS) Root() (Checksum, error) {
	rootObj, err := vfs.store.Get(rootKey)
	if err != nil {
		return Checksum{}, err
	}
	rootFile, ok := rootObj.(File)
	if !ok {
		return Checksum{}, fmt.Errorf("root object not of file type")
	}
	root := []byte((string)(rootFile))
	if len(root) != 32 {
		return Checksum{}, fmt.Errorf("root key unexpected length")
	}
	return (Checksum)(root), nil
}

func (vfs *VFS) setRoot(id Checksum) (Checksum, error) {
	idFile := (File)(string(id[:]))
	err := vfs.store.set(rootKey, &idFile)
	if err != nil {
		return Checksum{}, err
	}
	return rootKey, nil
}

func (vfs *VFS) Branch(id Checksum) (*Branch, error) {
	obj, err := vfs.store.Get(id)
	if err != nil {
		return nil, err
	}
	branch, ok := obj.(Branch)
	if !ok {
		return nil, fmt.Errorf("object at location '%x' is not a branch", id)
	}
	return &branch, nil
}

func (vfs *VFS) File(id Checksum) (*File, error) {
	obj, err := vfs.store.Get(id)
	if err != nil {
		return nil, err
	}
	file, ok := obj.(File)
	if !ok {
		return nil, fmt.Errorf("object at location '%x' is not a file", id)
	}
	return &file, nil
}

func (vfs *VFS) Get(path *Path) (*File, error) {
	curr, err := vfs.Root()
	if err != nil {
		return nil, err
	}

	for _, name := range *path {
		for j := 0; j < 32; j++ {
			b, err := vfs.Branch(curr)
			if err != nil {
				return nil, err
			}
			p := b.next(name.index(j))
			if p == nil {
				// we've reached a terminal node
				return nil, fmt.Errorf("directory does not exist")
			}
			if p.name != nil {
				if name.equals(p.name) {
					break
				} else {
					return nil, fmt.Errorf("directory does not exist")
				}
			}
			curr = p.id
		}
	}

	return vfs.File(curr)
}

type routeItem struct {
	// the index in the path to iterate to
	route *route
	pathIndex int
	branch *Branch
	index byte
	dirName *Name
}

func (r routeItem) pathKey() string {
	return r.route.pathString[:r.pathIndex]
}

type route struct {
	q []routeItem
	vfs *VFS
	pathString string
	path *Path
	last *routeItem
}

func initRoute(vfs *VFS, path *Path) *route {
	r := route{}
	r.vfs = vfs
	r.pathString = path.encoded()
	r.path = path
	r.last = nil
	r.q = make([]routeItem, len(r.pathString))
	for i := 0; i < len(r.q); i++ {
		item := routeItem{}
		item.route = &r
		item.pathIndex = i
		item.branch = nil
		item.index = r.pathString[i]
		r.q[i] = item
	}
	return &r
}

func (r route) isEmpty() bool {
	return len(r.q) <= 0
}

func (r *route) pop() (*routeItem, error) {
	old := r.q
	n := len(old)
	last := old[n-1]
	r.q = old[0 : n-1]
	r.last = &last
	return &last, nil
}

func (vfs VFS) route(path *Path) (*route, error) {
	route := initRoute(&vfs, path)
	key, err := vfs.Root()
	if err != nil {
		return nil, err
	}
	skip := false
	for i, item := range route.q {
		if skip && i % 64 > 0 {
			continue
		} else if i % 64 == 0 {
			skip = true
		}
		item.branch, err = vfs.Branch(key)
		if err != nil {
			return nil, err
		}
		bKey := route.pathString[i]
		bn := item.branch[bKey]
		if bn == nil {
			skip = true
		} else if bn.name != nil {
			pathName := (*path)[i / 16]
			if *bn.name == pathName {
				skip = true
			} else {
			}
		}
		key = bn.id
	}
	return route, nil
}

func (vfs *VFS) Gets(paths []*Path) ([]*File, error) {
	// TODO
	return nil, nil
}

func (vfs *VFS) Set(path *Path, file *File) (Checksum, error) {
	// then generate the planned route for updating
	route, err := vfs.route(path)

	// start by claiming all subpaths
	id := vfs.cache.Claim(route.pathString)

	route.populate()

	if err != nil {
		return Checksum{}, err
	}

	nextKey, err := vfs.store.Put(file)
	if err != nil {
		return Checksum{}, err
	}

	// iterate from the tail of the tree to the start
	for !route.isEmpty() {
		item, err := route.pop()
		if err != nil {
			return Checksum{}, err
		}

		key := item.pathKey()
		// Refresh the route if someone else beat us to it
		if newB := vfs.cache.Acquire(key, id); newB != nil {
			item.branch = newB
		}

		// Skip if our key doesn't exist yet
		if item.branch == nil {
			continue
		}

		// Generate the new branch
		b := item.branch
		b[item.index] = &branchNode{nextKey,item.dirName}
		nextKey = b.checksum()

		if item.pathIndex == 0 {
			vfs.setRoot(nextKey)
		} else {
			// Save the object to our store
			vfs.store.Put(b)
		}

		// Release lock on current path
		vfs.cache.Release(key, b)
	}

	return nextKey, nil
}

// func (vfs *VFS) Sets(paths []*Path, files []*File) (Id, error) {
// 	// BIG TODO
// 	// this is a whole lot harder...
// 	return Id{}, nil
// }
