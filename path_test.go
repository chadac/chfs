package chfs

import (
	"testing"
)

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
