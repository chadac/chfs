package chfs

type FileObj struct {
	path Path

	file File
	executable bool
}

type ChFS interface {
	// create an empty file system
	Init() (Ref, error)

	// primitive accessors
	GetRef(name string) (Ref, error)
	GetIndex(id Checksum) (Index, error)
	GetTree(id Checksum) (Tree, error)
	GetFile(id Checksum) (File, error)

	// Creates or updates an existing index
	UpdateRef(name string, id Checksum) (Ref, error)

	// Returns the main "root" of the tree, used for standard updates
	Head() (Ref, error)

	// An update is an atomic operation that changes the structure of a
	// tree in multiple places simultaneously
	Write(ref string, files []FileObj) (Checksum, error)

	// Read file(s)
	Read(ref string, paths []Path) ([]FileObj, error)

	// List a directory
	ListDir(ref string, paths []Path, recursive bool) ([][]FileObj, error)

	// // Merge replays the changes from `source` onto the given `ref`
	// Merge(ref string, source Index) (Index, error)
}
