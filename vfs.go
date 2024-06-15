package vfs

// // The root is stored at a single predictable location
// // It is the one mutable part of the store. This is used to track the root of
// // the trie for consistency.
// var rootKey = Checksum{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// type VersionFS struct {
// 	Index Store[Index]
// 	Branch Store[Branch]
// 	File Store[File]
// 	EventLog Log[WriteEvent]
// }

// type SFS interface {
// 	// set the vfs to an empty tree
// 	Ref(name string) (*Checksum, error)
// 	Reset() (*Checksum, error)
// 	setRoot(newRoot *Branch) (*Checksum, error)
// 	RootKey() (*Checksum, error)
// 	Root() (*Branch, error)
// 	// List(path *Path, recursive bool) ([]*Entry, error)
// 	Read(path *Path) (*Checksum, error)
// 	Write(event Event) (*Checksum, error)

// 	GetTree(paths []*Path) *Tree
// 	route(path *Path)
// }

// type Entry struct {
// 	path *Path
// 	file *File
// }

// // func vfsBranch(store store, id *Checksum) (*Branch, error) {
// // }

// // func vfsFile(store store, id *Checksum) (*File, error) {
// // }

// // func vfsGet(store Store, path *Path) (*File, error) {
// // }
