package fileutil_test

import (
	"testing"

	. "github.com/petergtz/pegomock"
	"github.com/relnod/fsa"
	"github.com/relnod/fsa/testutil"
	"github.com/stretchr/testify/assert"

	"github.com/relnod/dotm/pkg/fileutil"
	"github.com/relnod/dotm/pkg/mock"
)

func TestRecTraverseDir(t *testing.T) {
	RegisterMockTestingT(t)

	var tests = []struct {
		desc         string
		files        string
		visitorCalls [][]string
	}{
		{
			"No Visit calls for empty directories",
			`
				a/
				b/
			`,
			nil,
		},
		{
			"Simple Visit call",
			"a/b",
			[][]string{
				{"a", "b"},
			},
		},
		{
			"Multiple Visit calls in nested directories",
			`
				a/a
				a/b/c
				b/a/s/d
			`,
			[][]string{
				{"a", "a"},
				{"a/b", "c"},
				{"b/a/s", "d"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			fs := fsa.NewTempFs(fsa.NewOsFs())
			defer fs.Cleanup()

			err := testutil.CreateFiles(fs, test.files)
			assert.NoError(tt, err)

			visitor := mock.NewMockVisitor()

			err = fileutil.RecTraverseDir(fs, "", "", visitor)
			assert.NoError(tt, err)

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
		from string
		to   string
		dry  bool
	}{
		{
			"Simple linking works",
			"/a",
			"/b",
			false,
		},
		{
			"It doesn't link when doing a dry run",
			"/d",
			"/e",
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			fs := fsa.NewTempFs(fsa.NewOsFs())
			defer fs.Cleanup()

			assert.NoError(tt, testutil.CreateFiles(fs, test.from))
			assert.NoError(tt, fileutil.Link(fs, test.from, test.to, test.dry))
			assert.Equal(tt, testutil.IsSymlink(fs, test.to), !test.dry)
		})
	}
}

func TestUnlink(t *testing.T) {
	var tests = []struct {
		desc string
		file string
		dry  bool
	}{
		{
			"Simple unlinking works",
			"/a",
			false,
		},
		{
			"It doesn't unlink when doing a dry run",
			"/b",
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			fs := fsa.NewTempFs(fsa.NewOsFs())
			defer fs.Cleanup()

			assert.NoError(tt, testutil.CreateFiles(fs, test.file))
			assert.NoError(tt, fileutil.Unlink(fs, test.file, test.dry))
			assert.Equal(tt, testutil.FileExists(fs, test.file), test.dry)
		})
	}
}

func TestBackup(t *testing.T) {
	var tests = []struct {
		desc       string
		file       string
		dry        bool
		filesAfter string
	}{
		{
			"Simple backup works",
			"/a",
			false,
			`/a:deleted
			/a.backup`,
		},
		{
			"It doesn't backup when doing a dry run",
			"/b",
			true,
			"/b",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			fs := fsa.NewTempFs(fsa.NewOsFs())
			defer fs.Cleanup()

			assert.NoError(tt, testutil.CreateFiles(fs, test.file))
			assert.NoError(tt, fileutil.Backup(fs, test.file, test.dry))
			assert.NoError(tt, testutil.CheckFiles(fs, test.filesAfter))
		})
	}
}

func TestRestoreBackup(t *testing.T) {
	var tests = []struct {
		desc       string
		file       string
		dry        bool
		filesAfter string
	}{
		{
			"Restoring a backup works",
			"/a.backup",
			false,
			`/a.backup:deleted
			/a`,
		},
		{
			"It doesn't restore when doing a dry run",
			"/a.backup",
			true,
			"/a.backup",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			fs := fsa.NewTempFs(fsa.NewOsFs())
			defer fs.Cleanup()

			assert.NoError(tt, testutil.CreateFiles(fs, test.file))
			assert.NoError(tt, fileutil.RestoreBackup(fs, "/a", test.dry))
			assert.NoError(tt, testutil.CheckFiles(fs, test.filesAfter))
		})
	}
}
