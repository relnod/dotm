package profile_test

import (
	"path/filepath"
	"testing"

	"github.com/relnod/fsa"
	"github.com/relnod/fsa/testutil"
	"github.com/stretchr/testify/assert"

	"github.com/relnod/dotm/pkg/profile"
)

func TestActionLink(t *testing.T) {
	var tests = []struct {
		desc   string
		files  string
		source string
		dest   string
		name   string
	}{
		{
			"Can create simple symlink",
			"a/b",
			"a",
			"",
			"b",
		},
		{
			"Can create nested symlink",
			"foo/bar/blub",
			"foo/bar/",
			"foo/barnew/",
			"blub",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			fs := fsa.NewTempFs(fsa.NewOsFs())
			defer fs.Cleanup()

			assert.NoError(tt, testutil.CreateFiles(fs, test.files))

			action := profile.NewLinkAction(fs, nil)
			assert.NoError(tt, action.Run(
				test.source,
				test.dest,
				test.name,
			))
			assert.True(tt, testutil.IsSymlink(fs, filepath.Join(test.dest, test.name)))
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
	}{
		{
			"Can delete simple symlink",
			"a/b",
			"a",
			"",
			"b",
		},
		{
			"Can delete nested symlink",
			"foo/bar/blub:ln",
			"foo/bar/",
			"foo/bar/",
			"blub",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			fs := fsa.NewTempFs(fsa.NewOsFs())
			defer fs.Cleanup()

			assert.NoError(tt, testutil.CreateFiles(fs, test.files))

			action := profile.NewUnlinkAction(fs, nil)
			assert.NoError(tt, action.Run(
				test.source,
				test.dest,
				test.name,
			))
			assert.False(tt, testutil.FileExists(fs, filepath.Join(test.dest, test.name)))
		})
	}
}
