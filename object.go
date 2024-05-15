package vfs

// object is what is used to index a store
//
// since I have a bunch of stuff that ends up being key-value stores
// indexed by a checksum, might as well just make that a data
// structure
type Object interface {
	key() *Checksum
}
