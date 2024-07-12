package chfs

import (
	"fmt"
	"strings"
	"testing"
	"math/rand"
)

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
	st := NewSubTree(paths, false)
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
	st2 := NewSubTree(paths, true)
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

func BenchmarkSingleNewSubTree(b *testing.B) {
	paths := []Path{
		NewPath("/a/b/c"),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewSubTree(paths, false)
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomName(n int) string {
	b := make([]rune, n)
	for i := 0; i < n; i++ {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// TODO: This is slow!
// Let's make it faster
func BenchmarkSubTreeManyFilesOneDir(b *testing.B) {
	numFiles := 1000
	paths := make([]Path, numFiles)
	for i := 0; i < numFiles; i++ {
		paths[i] = NewPath(fmt.Sprintf("/a/b/c/d/%s", randomName(20)))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewSubTree(paths, true)
	}
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewSubTree(paths, true)
	}
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewSubTree(paths, false)
	}
}
