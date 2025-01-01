package chfs

import (
)

type Value[K comparable] interface {
	Key() K
}

/**
 * A key-value store that allows for bulk inserts/gets
 **/
type Store[K comparable, V Value[K]] interface {
	Get(key K) (*V, error)
	Gets(keys []K) ([]*V, error)
	Put(value V) (*K, error)
	Puts(values []V) ([]*K, error)
}


// stores the tree data structure, critical to chfs
type TreeStore Store[Checksum, Tree]

// stores indices, which are references to subpaths
type IndexStore Store[Checksum, Index]

// stores files
type FileStore Store[Checksum, File]

// stores references such as `HEAD` and such
type RefStore Store[string, Ref]


// todo: make this parallel
func listDir(store TreeStore, root Checksum, prefix Path, recursive bool) ([]PathObject, error) {
	tree, err := store.Get(root)
	if err != nil {
		return nil, err
	}

	// this is deterministic right now... let's preserve that behavior
	result := make([]PathObject, 0)
	for _, b := range tree.b {
		if b == nil {
			continue
		}
		// fmt.Printf("%s %b\n", b.id.repr(), b.obj.objType)
		doRecurse := true
		if b.obj != nil {
			result = append(result, NewPathObject(prefix, b))
			doRecurse = b.obj.objType == DirType && recursive
		}
		// fmt.Println(doRecurse)
		if doRecurse {
			newItems, err := listDir(store, b.id, prefix.Append(b.obj.name), recursive)
			if err != nil {
				return nil, err
			}
			result = append(result, newItems...)
		}
	}

	return result, nil
}

func ListDir(store TreeStore, root Checksum, recursive bool) ([]PathObject, error) {
	return listDir(store, root, NewPath("/"), recursive)
}
