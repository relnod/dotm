package dotfiles_test

import (
	"path/filepath"
	"testing"

	"github.com/relnod/dotm/pkg/dotfiles"
	"github.com/relnod/fsa"
	"github.com/relnod/fsa/osfs"
	"github.com/relnod/fsa/tempfs"
	"github.com/relnod/fsa/testutil"
	"github.com/stretchr/testify/assert"
)

func TestActionLink(t *testing.T) {
	var tests = []struct {
		desc   string
		files  string
		source string
		dest   string
		name   string
		err    error
	}{
		{
			"Can create simple symlink",
			"a/b",
			"a",
			"",
			"b",
			nil,
		},
		{
			"Can create nested symlink",
			"foo/bar/blub",
			"foo/bar/",
			"foo/barnew/",
			"blub",
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			fs := tempfs.New(osfs.New())
			defer fs.Cleanup()

			err := fsa.CreateFiles(fs, test.files)
			assert.NoError(tt, err)

			action := dotfiles.NewLinkAction(fs, false)
			err = action.Run(
				test.source,
				test.dest,
				test.name,
			)

			assert.Equal(tt, test.err, err)

			testutil.IsSymlink(tt, fs, err == nil, filepath.Join(test.dest, test.name))
		})
	}
}

func TestActionUnlink(t *testing.T) {
	var tests = []struct {
		desc   string
		files  string
		source string
		dest   string
		name   string
		err    error
	}{
		{
			"Can delete simple symlink",
			"a/b",
			"a",
			"",
			"b",
			nil,
		},
		{
			"Can delete nested symlink",
			"foo/bar/blub:ln",
			"foo/bar/",
			"foo/bar/",
			"blub",
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			fs := tempfs.New(osfs.New())
			defer fs.Cleanup()

			err := fsa.CreateFiles(fs, test.files)
			assert.NoError(tt, err)

			action := dotfiles.NewUnlinkAction(fs, false)
			err = action.Run(
				test.source,
				test.dest,
				test.name,
			)

			assert.Equal(tt, test.err, err)

			testutil.FileExists(tt, fs, err != nil, filepath.Join(test.dest, test.name))
		})
	}

}
