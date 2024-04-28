package vfs

import (
	"encoding/gob"
)

func init() {
	gob.Register(branch{})
	gob.Register(file{})
}
