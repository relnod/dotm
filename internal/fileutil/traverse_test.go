package fileutil_test

import (
	"testing"

	"github.com/relnod/dotm/internal/fileutil"
)

type visit struct {
	path string
	name string
}

type testVisitor struct {
	visits []visit
}

func (v *testVisitor) Visit(path, name string) error {
	v.visits = append(v.visits, visit{path, name})
	return nil
}

func TestRecTraverseDir(t *testing.T) {
	expectedVisits := []visit{
		visit{"testdata/a/a", "a/a"},
		visit{"testdata/a/b", "a/b"},
		visit{"testdata/c/a/a", "c/a/a"},
	}
	visitor := &testVisitor{
		visits: []visit{},
	}
	err := fileutil.RecTraverseDir("./testdata", visitor, "_")
	if err != nil {
		t.Fatal(err)
	}

	if len(expectedVisits) != len(visitor.visits) {
		t.Fatalf("expected %d visits. got %d.\nVisits:\n%v", len(expectedVisits), len(visitor.visits), visitor.visits)
	}
	for i, v := range visitor.visits {
		if v.path != expectedVisits[i].path || v.name != expectedVisits[i].name {
			t.Fatalf("expected visit %d to be %v. got %v", i, expectedVisits[i], v)
		}
	}
}
