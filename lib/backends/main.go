package backends

import (
	"github.com/chadac/vfs/lib"
)

type backend interface {
	put(key checksum, value []byte) bool
	get(key checksum) []byte
}
