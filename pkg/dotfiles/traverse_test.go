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
		excluded    []string
		actionCalls [][]string
	}{
		{
			"No action calls for empty directories",
			`a/
			 b/`,
			nil,
			nil,
		},
		{
			"Simple action call",
			"a/a",
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
			[][]string{
				{"b", "", "d"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			fs := fsa.NewTempFs(fsa.NewOsFs())
			defer fs.Cleanup()

			err := testutil.CreateFiles(fs, test.files)
			assert.NoError(tt, err)

			action := mock.NewMockAction()
			traverser := dotfiles.NewTraverser(fs, test.excluded)
			err = traverser.Traverse("", "", action)
			assert.NoError(tt, err)

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
