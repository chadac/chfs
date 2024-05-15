package vfs

// The root is stored at a single predictable location
// It is the one mutable part of the store. This is used to track the root of
// the trie for consistency.
var rootKey = Checksum{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

type VFS interface {
	// set the vfs to an empty tree
	Ref(name string) (*Checksum, error)
	Reset() (*Checksum, error)
	setRoot(newRoot *Branch) (*Checksum, error)
	RootKey() (*Checksum, error)
	Root() (*Branch, error)
	// List(path *Path, recursive bool)
	Read(path *Path) (*Checksum, error)
	// Gets(paths []*Path) ([]*File, error)
	Write(path *Path, checksum *Checksum) (*Checksum, error)
	// Sets(paths []*Path, files []*File) ([]*Checksum, error)
}

// func vfsBranch(store store, id *Checksum) (*Branch, error) {
// }

// func vfsFile(store store, id *Checksum) (*File, error) {
// }

// func vfsGet(store Store, path *Path) (*File, error) {
// }
