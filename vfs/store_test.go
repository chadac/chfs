package vfs

import (
	// "fmt"
	"testing"
)

func TestInMemoryStoreGet(t *testing.T) {
	s := NewInMemoryStore()
	var expected File = "test"
	key := expected.checksum()
	s.mem[key] = expected
	actual, err := s.Get(expected.checksum())
	if err != nil {
		t.Fatalf(`error: %s`, err)
	}
	if expected != actual {
		t.Fatalf(`expected '%s', got '%s'`, expected, actual)
	}
}

func TestInMemoryStorePut(t *testing.T) {
	s := NewInMemoryStore()
	var expected File = "test"
	s.Put(expected)
	s.Put((File)("junk"))
	s.Put((File)("extra"))
	s.Put((File)("fake-data"))
	actual, ok := s.mem[expected.checksum()]
	if !ok {
		t.Fatalf(`could not find '%s' in store`, expected)
	} else if expected != actual {
		t.Fatalf(`expected '%s', got '%s'`, expected, actual)
	}
}
