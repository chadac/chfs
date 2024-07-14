package chfs

import (
	"fmt"
	"strings"
	"testing"
	"math/rand"
)

// this tests just subtree logic, so no extra stuff here
func NewBool() bool { return true }

func nodeConsistencyTest[T any](n *PNode[T], t *testing.T) {
	// test character consistency
	if !n.IsRoot() {
		if (*n.name).Index(n.nameIndex) != n.char {
			t.Fatalf("at %s: turns out storing unnecessary crap in memory causes issues! put .char behind a function you idiot", n)
		}
	}

	// test pathIndex/nameIndex consistency
	if n.pathIndex > 0 {
		if n.prevDir == nil {
			// this mean
			t.Fatalf("at %s: prevDir is nil", n)
		}
		if n.prevTail == nil {
			t.Fatalf("as %s: prevTail is nil", n)
		}
		if n.prevDir.pathIndex != n.pathIndex - 1 || n.prevDir.nameIndex != 0 {
			t.Fatalf("at %s: prevDir %s malformed indices", n, n.prevDir)
		}
		if n.prevTail.pathIndex != n.pathIndex - 1 || n.prevTail.nameIndex != NameSize-1 {
			t.Fatalf("at %s: prevTail %s malformed indices", n, n.prevTail)
		}
	}
	if n.nextDir != nil {
		if n.nextDir.pathIndex != n.pathIndex + 1 || n.nextDir.nameIndex != 0 {
			t.Fatalf("at %s: nextDir %s malformed indices", n, n.nextDir)
		}
	}

	if len(n.next) <= 0 && n.nameIndex != NameSize - 1 {
		t.Fatalf("at %s: tree terminates before it fully populated its name", n)
	}

	// test all subsequent nodes of the tree
	for i, next := range n.next {
		if next == nil {
			t.Fatalf("%s has nil node at index %d", n, i)
		}
		if next.prev != n {
			t.Fatalf("at %s.next: next(%s).prev != n", n, n.next)
		}
		nextPathIndex, nextNameIndex := n.nextIndices()
		if next.pathIndex != nextPathIndex || nextNameIndex != next.nameIndex {
			t.Fatalf("at %s.next(%s): expected indices (%d, %d), got (%d, %d)", n, next, nextPathIndex, nextNameIndex, next.pathIndex, next.nameIndex)
		}
		nodeConsistencyTest(next, t)
	}
}

func TestSubTreeConsistency(t *testing.T) {
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

	st := NewSubTree(paths, NewBool)

	if len(st.Leafs()) != len(paths) {
		t.Fatalf("error: expected %d leafs, got %d", len(st.Leafs()), len(paths))
	}

	leafFailed := false
	for i, leaf := range st.Leafs() {
		if leaf == nil {
			leafFailed = true
			fmt.Printf("path at position %d is missing a leaf\n", i)
		}
	}
	if leafFailed {
		t.Error("tree inconsistent; at least one leaf is nil\n")
	}

	nodeConsistencyTest(st.root, t)
}

func TestNewSubTree(t *testing.T) {
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

	failed := false
	st := NewSubTree(paths, NewBool)
	for i, leaf := range st.Leafs() {
		if leaf == nil {
			failed = true
			fmt.Printf("path at position %d is missing a leaf\n", i)
		}
	}
	if failed {
		t.Error("failed to create tree\n")
	}

	failed = false
	st2 := NewSubTree(paths, NewBool)
	for i, leaf := range st2.Leafs() {
		if leaf == nil {
			failed = true
			fmt.Printf("path at position %d is missing a leaf\n", i)
		}
	}
	if failed {
		t.Error("failed to create tree with parallel\n")
	}
}

func doBenchCreate(paths []Path, b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewSubTree(paths, NewBool)
	}
}

func BenchmarkSingleNewSubTree(b *testing.B) {
	paths := []Path{
		NewPath("/a/b/c"),
	}
	doBenchCreate(paths, b)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomName(n int) string {
	b := make([]rune, n)
	for i := 0; i < n; i++ {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// func BenchmarkBaseAllocation(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		dirCount := 5
// 		x := make([][]*PNode[int], dirCount)
// 		for d := 0; d < dirCount; d++ {
// 			x[d] = make([]*PNode[int], NameSize)
// 			for j := 0; j < NameSize; j++ {
// 				x[d][j] = &PNode[int]{}
// 			}
// 		}
// 	}
// }

// TODO: This is slow! Takes ~50ms to run, when it should be <1ms
// probably should preallocate nodes to make it faster
func BenchmarkSubTreeManyFilesOneDir(b *testing.B) {
	numFiles := 1000
	paths := make([]Path, numFiles)
	for i := 0; i < numFiles; i++ {
		paths[i] = NewPath(fmt.Sprintf("/a/b/c/d/%s", randomName(20)))
	}
	doBenchCreate(paths, b)
}

func BenchmarkNewSubTreeManyFiles(b *testing.B) {
	numFiles := 1000
	numDirs := 5
	paths := make([]Path, numFiles)

	for i := 0; i < numFiles; i++ {
		sb := strings.Builder{}
		for d := 0; d < numDirs; d++ {
			sb.WriteString("/")
			sb.WriteString(randomName(2))
		}
		paths[i] = NewPath(sb.String())
	}

	doBenchCreate(paths, b)
}

func BenchmarkNewSubTree(b *testing.B) {
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

	doBenchCreate(paths, b)
}
