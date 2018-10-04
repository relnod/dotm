package dotfiles_test

import (
	"testing"

	. "github.com/petergtz/pegomock"
	"github.com/relnod/fsa"
	"github.com/relnod/fsa/testutil"
	"github.com/stretchr/testify/assert"

	"github.com/relnod/dotm/pkg/dotfiles"
	"github.com/relnod/dotm/pkg/mock"
)

func TestTraverse(t *testing.T) {
	RegisterMockTestingT(t)

	var tests = []struct {
		desc        string
		files       string
		excludes    []string
		includes    []string
		actionCalls [][]string
	}{
		{
			"No action calls for empty directories",
			`a/
			 b/`,
			nil,
			nil,
			nil,
		},
		{
			"Simple action call",
			"a/a",
			nil,
			nil,
			[][]string{
				{"a", "", "a"},
			},
		},
		{
			"Multiple action calls in nested directories",
			`a/a
			 a/b/c
			 b/d`,
			nil,
			nil,
			[][]string{
				{"a", "", "a"},
				{"a/b", "b", "c"},
				{"b", "", "d"},
			},
		},
		{
			"Can exclude top level directories",
			`a/a,
			 a/b/c
			 b/d`,
			[]string{"a"},
			nil,
			[][]string{
				{"b", "", "d"},
			},
		},
		// FIXME: Those tests should pass, but, combined with the one above the
		// don't for some reason.
		// {
		// 	"Can include top level directories",
		// 	`a/a,
		// 	 a/b/c
		// 	 b/d`,
		// 	nil,
		// 	[]string{"a"},
		// 	[][]string{
		// 		{"a", "", "a"},
		// 		{"a/b", "b", "c"},
		// 	},
		// },
		// {
		// 	"Includes and excludes",
		// 	`a/a,
		// 	 b/b/c
		// 	 c/d`,
		// 	[]string{"a"},
		// 	[]string{"a", "b"},
		// 	[][]string{
		// 		{"a", "", "a"},
		// 	},
		// },
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			fs := fsa.NewTempFs(fsa.NewOsFs())
			defer fs.Cleanup()

			assert.NoError(tt, testutil.CreateFiles(fs, test.files))

			action := mock.NewMockAction()
			traverser := dotfiles.NewTraverser(fs, test.excludes, test.includes)
			assert.NoError(tt, traverser.Traverse("", "", action))

			action.VerifyWasCalled(Times(len(test.actionCalls))).Run(AnyString(), AnyString(), AnyString())

			if test.actionCalls != nil {
				inOrderContext := new(InOrderContext)
				for _, call := range test.actionCalls {
					action.VerifyWasCalledInOrder(Once(), inOrderContext).Run(
						call[0],
						call[1],
						call[2],
					)
				}
			}
		})
	}
}
