package vfs

import (
	"fmt"
	"sync"
)

type Registry struct {
	Backend *Backend
	writeCount uint64
	writes sync.Map
	writeMu sync.Mutex
}

type FileNotFoundError struct {
	// used internally to determine if the branch conflicts
	isBranchConflict bool
}
func (e FileNotFoundError) Error() string {
	return "file not found"
}

func New(backend *Backend) *Registry {
	r := new(Registry)
	r.Backend = backend
	return r
}

func (r *Registry) branch(id *id) (*branch, error) {
	return nil, nil
	// var data, err = r.backend.get(id)
	// if err != nil {
	// 	return branch{}, err
	// }
	// if len(data) != baseSize * (chkSize + 1) {
	// 	return branch{}, fmt.Errorf("branch '%+v' has invalid size", id)
	// }
	// var refs [baseSize]ref
	// for i := 0; i < baseSize; i++ {
	// 	var offset = i * (chkSize + 1)
	// 	refs[i] = ref{data[offset], (checksum)(data[offset+1:offset+chkSize+1])}
	// }
	// return refs, nil
}

func (r *Registry) file(id *id) (*file, error) {
	return nil, nil
	// return r.backend.get(id)
}

// returns the ID of the root node
// interestingly, this is stored in the registry as a file. weird, eh?
func (r *Registry) root() (*id, error) {
	root, err := r.file(&rootKey)
	if err != nil {
		return nil, err
	}
	return root.hash, nil
}

func (r *Registry) dirRoute(root *id, dir *directory) (routeComponent, error) {
	curr := root
	ret := routeComponent{[]*id{root}, nil}
	for _, char := range dir.x {
		b, err := r.branch(curr)
		if err != nil {
			return ret, err
		}
		next, err := b.get(char)
		if err != nil {
			return ret, err
		}
		if next.name != nil {
			if next.name == dir.name {
				ret.next = next.ref
				break
			} else {
				return ret, FileNotFoundError{true}
			}
		}
		if next.ref == nil {
			return ret, FileNotFoundError{false}
		}
		ret.subpath = append(ret.subpath, next.ref)
		curr = next.ref
	}
	return ret, nil
}

func (r *Registry) route(root *id, p *path) (route, error) {
	curr := root
	ret := route{}
	for _, dir := range *p {
		routePart, err := r.dirRoute(curr, &dir)
		ret = append(ret, routePart)
		if err != nil {
			return ret, err
		}
		curr = routePart.next
	}
	return ret, nil
}

func (r *Registry) get(root *id, path *path) (*file, error) {
	rt, err := r.route(root, path)
	if err != nil {
		return nil, err
	}
	id := rt.file()
	if id == nil {
		return nil, fmt.Errorf("path does not return a valid file")
	}
	return r.file(id)
}

func (r *Registry) del(root *id, path *path) error {
	return nil
}

func (r *Registry) set(root *id, path *path, file *file) (*id, error) {
	rt, err := r.route(root, path)
	if fnf, ok := err.(FileNotFoundError); ok {
		// create the file
		if fnf.isBranchConflict {
			// this means that we reached a split... so we gotta create a
			// new branch with the proper split
		} else {
		}
	} else if err != nil {
		return nil, err
	}
	// incrementally replace the file
	return nil, nil
}

// // update the registry with the given data at the given key
// func (r *Registry) set(root *id, path []ref, data []byte) (checksum, error) {
// 	// TODO: add locking mechanism for concurrency
// 	p, err := r.path(id, path)
// 	if err != nil {
// 		return checksum{}, err
// 	}
// 	next, err := r.backend.put(data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for i := len(p)-2; i >= 0; i-- {
// 		b := r.branch(p[i])
// 	}
// 	return checksum{}, nil
// }

// func (r *Registry) bulkGet(root checksum, path []ref) ([]checksum, error) {
// 	return nil, nil
// }

// func (r *Registry) bulkSet(root checksum, path []ref, data []byte) (checksum, error) {
// 	return checksum{}, nil
// }

// func (r *Registry) list(root checksum, prefix []ref) ([]*node, error) {
// 	return nil, nil
// }

// func main() {
// }
