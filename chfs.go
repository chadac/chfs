package chfs

type PathObject struct {
	directory Path
	filename Name
	fileId Checksum

	executable bool
}

func NewPathObject(directory Path, branch *Branch) PathObject {
	return PathObject{
		directory,
		branch.obj.name,
		branch.id,
		branch.obj.executable,
	}
}

func (p PathObject) Path() Path {
	return p.directory.Append(p.filename)
}

func (p PathObject) Branch() *Branch {
	file := NewFile(p.filename)
	file.executable = p.executable
	branch := Branch{
		p.fileId,
		&file,
	}
	return &branch
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
	Write(ref string, files []PathObject) (Checksum, error)

	// Read file(s)
	Read(ref string, paths []Path) ([]PathObject, error)

	// List a directory
	ListDir(ref string, paths []Path, recursive bool) ([][]PathObject, error)

	// // Merge replays the changes from `source` onto the given `ref`
	// Merge(ref string, source Index) (Index, error)
}

// Creates an empty tree
func Init(ref RefStore, tree TreeStore) error {
	root := EmptyTree()
	id, err := tree.Put(*root)
	if err != nil {
		return err
	}
	rootRef := Ref{"HEAD",*id}
	ref.Put(rootRef)
	return nil
}
