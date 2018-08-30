package dotfiles_test

import (
	"testing"

	"github.com/relnod/dotm/internal/testutil"
	"github.com/relnod/dotm/internal/testutil/assert"
	"github.com/relnod/dotm/pkg/dotfiles"
)

func TestActionLink(t *testing.T) {
	var tests = []struct {
		desc          string
		fileStructure testutil.FileStructure
		source        string
		dest          string
		name          string
		err           error
	}{
		{
			"Can create simple symlink",
			testutil.FileStructure{
				"a/b",
			},
			"a",
			"",
			"b",
			nil,
		},
		{
			"Can create nested symlink",
			testutil.FileStructure{
				"foo/bar/blub",
			},
			"foo/bar/",
			"foo/bar/",
			"blub",
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			source := testutil.NewFileSystem()
			dest := testutil.NewFileSystem()
			defer source.Cleanup()
			defer dest.Cleanup()

			source.CreateFromFileStructure(test.fileStructure)

			action := dotfiles.NewLinkAction(false)
			err := action.Run(
				source.Path(test.source),
				dest.Path(test.dest),
				test.name,
			)

			assert.ErrorEquals(tt, err, test.err)

			if test.err == nil {
				assert.IsSymlink(tt, dest.Join(test.dest, test.name))
			}
		})
	}
}

func TestActionUnlink(t *testing.T) {
	var tests = []struct {
		desc          string
		fileStructure testutil.FileStructure
		source        string
		dest          string
		name          string
		err           error
	}{
		{
			"Can delete simple symlink",
			testutil.FileStructure{
				"a/b",
			},
			"a",
			"",
			"b",
			nil,
		},
		{
			"Can delete nested symlink",
			testutil.FileStructure{
				"foo/bar/blub",
			},
			"foo/bar/",
			"foo/bar/",
			"blub",
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			source := testutil.NewFileSystem()
			dest := testutil.NewFileSystem()
			defer source.Cleanup()
			defer dest.Cleanup()

			source.CreateFromFileStructure(test.fileStructure)

			action := dotfiles.NewUnlinkAction(false)
			err := action.Run(
				source.Path(test.source),
				dest.Path(test.dest),
				test.name,
			)

			assert.ErrorEquals(tt, err, test.err)

			if test.err == nil {
				assert.PathNotExists(tt, dest.Join(test.dest, test.name))
			}
		})
	}

}
