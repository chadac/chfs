package vfs

import (
	// "errors"
)

type Backend interface {
	Put(value []byte) (id, error)
	Puts(value [][]byte) ([]id, error)
	Get(key id) ([]byte, error)
	Gets(keys []id) ([][]byte, error)
}
