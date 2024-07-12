package chfs

type Value[K comparable] interface {
	Key() K
}

/**
 * A key-value store that allows for bulk inserts/gets
 **/
type Store[K comparable, V Value[K]] interface {
	Get(key K) (V, error)
	Gets(keys []K) []V
	Put(value V) (K, error)
	Puts(values []V) []K
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
func listDir(store TreeStore, root Checksum, recursive bool) ([]*Branch, error) {
	tree, err := store.Get(root)
	if err != nil {
		return nil, err
	}

	// this is deterministic right now... let's preserve that behavior
	result := make([]*Branch, 0)
	for _, b := range tree.b {
		if b == nil {
			continue
		}
		doRecurse := true
		if b.obj != nil {
			result = append(result, b)
			doRecurse = b.obj.objType == DirType && recursive
		}
		if doRecurse {
			newItems, err := ListDir(store, b.id, recursive)
			if err != nil {
				return nil, err
			}
			result = append(result, newItems...)
		}
	}

	return result, nil
}

func ListDir(store TreeStore, root Checksum, recursive bool) ([]*Branch, error) {
	return listDir(store, root, recursive)
}

func Read(store TreeStore, root Checksum, paths []Path) ([]*Branch, error) {
}
