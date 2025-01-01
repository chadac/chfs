package chfs/store

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
