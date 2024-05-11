package vfs

// The root is stored at a single predictable location
// It is the one mutable part of the store. This is used to track the root of
// the trie for consistency.
var rootKey = Checksum{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

type VFS interface {
	// set the vfs to an empty tree
	Reset() (*Checksum, error)
	setRoot(newRoot *Branch) (*Checksum, error)
	RootKey() (*Checksum, error)
	Root() (*Branch, error)
	Get(path *Path) (*File, error)
	// Gets(paths []*Path) ([]*File, error)
	Set(path *Path, file *File) (*Checksum, error)
	// Sets(paths []*Path, files []*File) ([]*Checksum, error)
}

// func vfsBranch(store store, id *Checksum) (*Branch, error) {
// }

// func vfsFile(store store, id *Checksum) (*File, error) {
// }

// func vfsGet(store Store, path *Path) (*File, error) {
// }
