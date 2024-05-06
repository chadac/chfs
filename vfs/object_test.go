package vfs

import (
	"testing"
)

const sampleFile File = "{\"path\": \"https://example.org/file.txt\"}"
var sampleChecksum Checksum = [32]byte{91, 113, 23, 164, 237, 163, 132, 141, 89, 166, 82, 174, 177, 20, 34, 148, 147, 28, 249, 101, 35, 138, 16, 240, 58, 53, 216, 192, 246, 4, 3, 52};
const sampleRepr string = "5b7117a4eda3848d59a652aeb1142294931cf965238a10f03a35d8c0f6040334"


func TestChecksumEquals(t *testing.T) {
	if !sampleChecksum.equals(&sampleChecksum) {
		t.Fatalf(`'%x' does not equal itself`, sampleChecksum)
	}
}

func TestFileChecksum(t *testing.T) {
	if !sampleFile.checksum().equals(&sampleChecksum) {
		t.Fatalf(`expected: '%x', got '%x'`, sampleChecksum, sampleFile)
	}
}

func TestChecksumIndex(t *testing.T) {
	earr := []byte{5, 11, 7, 1, 1, 7, 10, 4, 14}
	for i, expected := range earr {
		actual := sampleChecksum.index(i)
		if actual != expected {
			t.Fatalf(`at index %d expected '%d' got '%d'`, i, expected, actual)
		}
	}
}

func TestPointerRepr(t *testing.T) {
}
