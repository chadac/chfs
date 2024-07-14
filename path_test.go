package chfs

import (
	"testing"
)

func TestPathEmpty(t *testing.T) {
	p := NewPath("/")
	if len(p) != 0 {
		t.Fatalf("what the fuck %d", len(p))
	}
}

func TestPathCreate(t *testing.T) {
	p := NewPath("/a/b/c")
	if len(p) != 3 {
		t.Fatalf("expected len(p) = 3, got '%d'", len(p))
	}
	expected := []string{"a", "b", "c"}
	for i, e := range expected {
		if p[i].raw != e {
			t.Fatalf("expected p[%d] = '%s', got %s", i, e, p[i].raw)
		}
	}
}

func TestPathString(t *testing.T) {
	expected := "/a/b/c"
	p1 := NewPath(expected)
	if p1.String() != expected {
		t.Fatalf("error: expected '%s', got '%s'", expected, p1.String())
	}
}
