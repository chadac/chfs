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
	paths PathGroup
}

func (t *Transaction) Execute() error {
	return nil
}
