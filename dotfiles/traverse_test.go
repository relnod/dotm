package dotfiles_test

import (
	"testing"

	. "github.com/petergtz/pegomock"

	"github.com/relnod/dotm/dotfiles"
	"github.com/relnod/dotm/internal/mock"
	"github.com/relnod/dotm/internal/testutil"
)

func TestTraverse(t *testing.T) {
	RegisterMockTestingT(t)

	var tests = []struct {
		desc          string
		fileStructure testutil.FileStructure
		excluded      []string
		actionCalls   [][]string
	}{
		{
			"No action calls for empty directories",
			testutil.FileStructure{
				"a/",
				"b/",
			},
			nil,
			nil,
		},
		{
			"Simple action call",
			testutil.FileStructure{
				"a/a",
			},
			nil,
			[][]string{
				{"a", "", "a"},
			},
		},
		{
			"Multiple action calls in nested directories",
			testutil.FileStructure{
				"a/a",
				"a/b/c",
				"b/d",
			},
			nil,
			[][]string{
				{"a", "", "a"},
				{"a/b", "b", "c"},
				{"b", "", "d"},
			},
		},
		{
			"Can exclude top level directories",
			testutil.FileStructure{
				"a/a",
				"a/b/c",
				"b/d",
			},
			[]string{"a"},
			[][]string{
				{"b", "", "d"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			source := testutil.NewFileSystem()
			dest := testutil.NewFileSystem()
			defer source.Cleanup()
			defer dest.Cleanup()

			source.CreateFromFileStructure(test.fileStructure)

			action := mock.NewMockAction()
			traverser := dotfiles.NewTraverser(test.excluded)
			traverser.Traverse(source.BasePath(), dest.BasePath(), action)

			action.VerifyWasCalled(Times(len(test.actionCalls))).Run(AnyString(), AnyString(), AnyString())

			if test.actionCalls != nil {
				inOrderContext := new(InOrderContext)
				for _, call := range test.actionCalls {
					action.VerifyWasCalledInOrder(Once(), inOrderContext).Run(
						source.Path(call[0]),
						dest.Path(call[1]),
						call[2],
					)
				}
			}

		})
	}
}
