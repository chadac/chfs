package main

import (
	// "errors"
	"fmt"
	// "io"
	// "net/http"
	// "os"
	"crypto/sha256"
)

const chkSize = sha256.Size
type checksum [chkSize]byte
var rootKey = [chkSize]byte{0}

func computeChecksum(data []byte) checksum {
	return checksum{}
}

// base of the encoding (i.e. base 16 by default)
const baseSize = 16
const branchSize = chkSize * baseSize;

type branch [baseSize]ref
type file []byte
type link struct {
	id checksum
	next checksum
	directory bool
}

type ref struct {
	kind byte
	id checksum
}

type registry struct {
	backend backend
}

type name [baseSize]byte;

func (b *branch) get(key byte) (ref, error) {
	if key > baseSize {
		return ref{}, fmt.Errorf("index '%+v' out of range", key)
	}
	return b[key], nil
}

func (r *registry) branch(id checksum) (branch, error) {
	var data, err = r.backend.get(id)
	if err != nil {
		return branch{}, err
	}
	if len(data) != baseSize * (chkSize + 1) {
		return branch{}, fmt.Errorf("branch '%+v' has invalid size", id)
	}
	var refs [baseSize]ref
	for i := 0; i < baseSize; i++ {
		var offset = i * (chkSize + 1)
		refs[i] = ref{data[offset], (checksum)(data[offset+1:offset+chkSize+1])}
	}
	return refs, nil
}

func (r *registry) file(id checksum) (file, error) {
	return r.backend.get(id)
}

func (r *registry) link(id checksum) (link, error) {
	data, err := r.backend.get(id)
	if err != nil {
		return link{}, err
	}
	if len(data) != 2 * chkSize {
		return link{}, fmt.Errorf("directory '%+v' has invalid size", id)
	}
	return link{(checksum)(data[0:chkSize]), (checksum)(data[chkSize:2*chkSize])}, nil
}

func (r *registry) root() (checksum, error) {
	d, err := r.backend.get(rootKey)
	if err != nil {
		return checksum{}, err
	}
	return (checksum)(d), nil
}

type branchRef struct {
	id checksum
	c byte
	value branch
}

func (r *registry) path(root checksum, path []name) (checksum, []branchRef, error) {
	p := []branchRef{}
	curr := root
	for idx, n := range path {
		for _, k := range n {
			b, err := r.branch(curr.id)
			if err != nil {
				return checksum{}, p, err
			}
			p = append(p, branchRef{curr, b})
			curr, err = b.get(k)
			if err != nil {
				return checksum{}, p, err
			}
			if curr.kind == dirType {
				break
			} else if curr.kind == fileType {
				if idx == len(path)-1 {
					return curr.id, p, nil
				} else {
					return checksum{}, p, fmt.Errorf("attempted to read file as directory")
				}
			} else if curr.kind != branchType {
				return checksum{}, p, fmt.Errorf("unexpected ref type '%+v'", curr.kind)
			}
		}
	}
	return checksum{}, p, fmt.Errorf("attempted to fetch node from directory")
}

func (r *registry) get(root checksum, path []name) (file, error) {
	id, _, err := r.path(root, path)
	if err != nil {
		return nil, err
	}
	return r.file(id)
}

// update the registry with the given data at the given key
func (r *registry) set(root checksum, path []ref, data []byte) (checksum, error) {
	// TODO: add locking mechanism for concurrency
	p, err := r.path(id, path)
	if err != nil {
		return checksum{}, err
	}
	next, err := r.backend.put(data)
	if err != nil {
		return nil, err
	}
	for i := len(p)-2; i >= 0; i-- {
		b := r.branch(p[i])
	}
	return checksum{}, nil
}

func (r *registry) bulkGet(root checksum, path []ref) ([]checksum, error) {
	return nil, nil
}

func (r *registry) bulkSet(root checksum, path []ref, data []byte) (checksum, error) {
	return checksum{}, nil
}

func (r *registry) list(root checksum, prefix []ref) ([]*node, error) {
	return nil, nil
}

func main() {
}
