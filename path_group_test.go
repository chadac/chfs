package chfs

import (
	"testing"
)

func TestPathGroupCreate(t *testing.T) {
	paths := []Path{
		NewPath("/a/b/c"),
		NewPath("/a/c/a"),
		NewPath("/a/c/b"),
		NewPath("/a/ba/b"),
		NewPath("/a/bb/b"),
		NewPath("/a/bc/b"),
		NewPath("/a/bd/b"),
		NewPath("/a/be/b"),
		NewPath("/a/bf/b"),
		NewPath("/a/bg/b"),
		NewPath("/a/bh/b"),
	}
	pg := CreatePathGroup(paths)
	// pg.root.print()
	if len(pg.Leafs()) != len(paths) {
		t.Fatalf("number of leafs does not match number of paths; expected %d, got %d", len(paths), len(pg.Leafs()))
	}
}

func BenchmarkPathGroupCreate(b *testing.B) {
	paths := []Path{
		NewPath("/a/b/c"),
		NewPath("/a/c/a"),
		NewPath("/a/c/b"),
		NewPath("/a/ba/b"),
		NewPath("/a/bb/b"),
		NewPath("/a/bc/b"),
		NewPath("/a/bd/b"),
		NewPath("/a/be/b"),
		NewPath("/a/bf/b"),
		NewPath("/a/bg/b"),
		NewPath("/a/bh/b"),
	}
	for i := 0; i < b.N; i++ {
		CreatePathGroup(paths)
	}
}
