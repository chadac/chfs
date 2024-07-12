package chfs

type RmAction struct {
	recursive bool
}

type CreateAction struct {
	file File
}

type Transaction struct {
	id int64
	ref Ref
	subtree SubTree
}

func (t *Transaction) Execute() error {
	return nil
}
