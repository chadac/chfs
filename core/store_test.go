package chfs

import (
	// "fmt"
	"testing"
)

func TestInMemoryStoreGet(t *testing.T) {
	s := NewInMemoryStore()
	key := sampleFile.checksum()
	s.mem[*key] = sampleFile
	actual, err := s.Get(key)
	if err != nil {
		t.Fatalf(`error: %s`, err)
	}
	if sampleFile != actual {
		t.Fatalf(`expected '%s', got '%s'`, sampleFile, actual)
	}
}

func TestInMemoryStorePut(t *testing.T) {
	s := NewInMemoryStore()
	var expected = FileFromString("test")
	s.Put(expected)
	s.Put(FileFromString("testing"))
	s.Put(FileFromString("extra"))
	s.Put(FileFromString("fake-data"))
	actual, ok := s.mem[*expected.checksum()]
	if !ok {
		t.Fatalf(`could not find '%s' in store`, expected)
	} else if expected != actual {
		t.Fatalf(`expected '%s', got '%s'`, expected, actual)
	}
}
