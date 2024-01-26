package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	// "github.com/minio/sha256-simd"
)

const chkSize = sha256.size
type checksum [chkSize]byte

// base of the encoding (i.e. base 16 by default)
const baseSize = 16

type node struct {
	id checksum
	data []byte
}

type branch struct {
	id checksum
	nodes [baseSize]checksum
}

type ref struct {
	id checksum
	name string
}

type backend interface {
	put(key checksum, value []byte) bool
	get(key checksum) []byte
}

type repo struct {
	backend *backend
}

func get(version string, path []ref) {
}

func put(root checksum, prefix []ref, data []byte) checksum {
}

func list(root checksum, prefix []ref) {
}


func main() {
}
