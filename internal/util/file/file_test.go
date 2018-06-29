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
			defer fs.Cleanup()

			fs.CreateFromFileStructure(test.fileStructure)

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

func TestLink(t *testing.T) {
	var tests = []struct {
		desc string
		fs   testutil.FileStructure
		from string
		to   string
		dry  bool
		err  error
	}{
		{
			"Simple linking works",
			testutil.FileStructure{"a"},
			"a",
			"b",
			false,
			nil,
		},
		{
			"It doesn't link when doing a dry run",
			testutil.FileStructure{"a"},
			"a",
			"b",
			true,
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			fs := testutil.NewFileSystem()
			defer fs.Cleanup()

			fs.CreateFromFileStructure(test.fs)

			destFS := testutil.NewFileSystem()
			defer destFS.Cleanup()

			err := file.Link(fs.Path(test.from), destFS.Path(test.to), test.dry)

			assert.ErrorEquals(tt, err, test.err)

			if err == nil && !test.dry {
				assert.IsSymlink(tt, destFS.Path(test.to))
			}
		})
	}
}
