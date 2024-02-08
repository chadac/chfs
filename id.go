package main

const chkSize = sha256.Size
type checksum [chkSize]byte
type id checksum
var rootKey = new([chkSize]byte)

// base of the encoding (i.e. base 16 by default)
const baseSize = 16

func computeChecksum(data []byte) checksum {
	return checksum{}
}
