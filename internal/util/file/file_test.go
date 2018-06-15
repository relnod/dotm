package file_test

import (
	"testing"

	. "github.com/petergtz/pegomock"

	"github.com/relnod/dotm/internal/mock"
	"github.com/relnod/dotm/internal/testutil"
	"github.com/relnod/dotm/internal/testutil/assert"
	"github.com/relnod/dotm/internal/util/file"
)

func TestRecTraverseDir(t *testing.T) {
	RegisterMockTestingT(t)

	var tests = []struct {
		desc          string
		fileStructure testutil.FileStructure
		visitorCalls  [][]string
	}{
		{
			"No Visit calls for empty directories",
			testutil.FileStructure{
				"a/",
				"b/",
			},
			nil,
		},
		{
			"Simple Visit call",
			testutil.FileStructure{
				"a/b",
			},
			[][]string{
				{"a", "b"},
			},
		},
		{
			"Multiple Visit calls in nested directories",
			testutil.FileStructure{
				"a/a",
				"a/b/c",
				"b/a/s/d",
			},
			[][]string{
				{"a", "a"},
				{"a/b", "c"},
				{"b/a/s", "d"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			fs := testutil.NewFileSystem()

			fs.CreateFromFileStructure(test.fileStructure)
			defer fs.Cleanup()

			visitor := mock.NewMockVisitor()

			err := file.RecTraverseDir(fs.BasePath(), "", visitor)
			assert.ErrorIsNil(tt, err)

			visitor.VerifyWasCalled(Times(len(test.visitorCalls))).Visit(AnyString(), AnyString())

			if test.visitorCalls != nil {
				inOrderContext := new(InOrderContext)
				for _, call := range test.visitorCalls {
					visitor.VerifyWasCalledInOrder(Once(), inOrderContext).Visit(
						call[0],
						call[1],
					)
				}
			}

		})
	}
}
