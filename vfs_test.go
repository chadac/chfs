package chfs

// import (
// 	"testing"
// )

// func newTestVFS() *VFS {
// 	store := NewInMemoryStore()
// 	cache := NewLocalWriteCache()
// 	vfs := VFS{store,cache}
// 	vfs.Reset()
// 	return &vfs
// }

// func TestVFSSetGetOnce(t *testing.T) {
// 	s := newTestVFS()
// 	path := NewPath("a/b/c")
// 	expected := FileFromString("test")
// 	_, err := s.Set(path, expected)
// 	if err != nil {
// 		t.Fatalf(`error: %s`, err)
// 	}
// 	actual, err := s.Get(path)
// 	if err != nil {
// 		t.Fatalf(`error: %s`, err)
// 	}
// 	if *expected != *actual {
// 		t.Fatalf(`expected '%s', got '%s'`, expected, *actual)
// 	}
// }
