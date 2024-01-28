package main

import (
	// "errors"
)

type backend interface {
	put(value []byte) (checksum, error)
	puts(value [][]byte) ([]checksum, error)
	get(key checksum) ([]byte, error)
	gets(keys []checksum) ([][]byte, error)
}
