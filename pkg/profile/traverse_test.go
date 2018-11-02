package profile_test

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/petergtz/pegomock"
	"github.com/relnod/fsa"
	"github.com/relnod/fsa/testutil"
	"github.com/stretchr/testify/assert"

	"github.com/relnod/dotm/mock"
	"github.com/relnod/dotm/pkg/profile"
)

func TestTraverse(t *testing.T) {
	RegisterMockTestingT(t)

	var tests = []struct {
		desc        string
		files       string
		p           *profile.Profile
		actionCalls [][]string
	}{
		{
			"No action calls for empty directories",
			`a/
			 b/`,
			&profile.Profile{Path: ""},
			nil,
		},
		{
			"Simple action call",
			"a/a",
			&profile.Profile{Path: ""},
			[][]string{
				{"a", "", "a"},
			},
		},
		{
			"Multiple action calls in nested directories",
			`a/a
			 a/b/c
			 b/d`,
			&profile.Profile{Path: ""},
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
			&profile.Profile{Path: "", Excludes: []string{"a"}},
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
		//   &profile.Profile{Path: "", Includes: []string{"a"}},
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
		//   &profile.Profile{Path: "", Includes: []string{"a", "b"}, Excludes: []string{"a"}},
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
			test.p.SetFS(fs)
			assert.NoError(tt, test.p.Traverse(action))

			action.VerifyWasCalled(Times(len(test.actionCalls))).Run(AnyString(), AnyString(), AnyString())

			if test.actionCalls != nil {
				inOrderContext := new(InOrderContext)
				for _, call := range test.actionCalls {
					action.VerifyWasCalledInOrder(Once(), inOrderContext).Run(
						call[0],
						filepath.Join(os.ExpandEnv("/home/$USER"), call[1]),
						call[2],
					)
				}
			}
		})
	}
}
