package fileutil_test

import (
	"testing"

	. "github.com/petergtz/pegomock"
	"github.com/relnod/fsa"
	"github.com/relnod/fsa/osfs"
	"github.com/relnod/fsa/tempfs"
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
			fs := tempfs.New(osfs.New())
			defer fs.Cleanup()

			err := fsa.CreateFiles(fs, test.files)
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
		err  error
	}{
		{
			"Simple linking works",
			"a",
			"b",
			false,
			nil,
		},
		{
			"It doesn't link when doing a dry run",
			"a",
			"b",
			true,
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			fs := tempfs.New(osfs.New())
			defer fs.Cleanup()

			_, err := fs.Create(test.from)
			assert.Equal(tt, test.err, err)

			err = fileutil.Link(fs, test.from, test.to, test.dry)
			assert.Equal(tt, test.err, err)
			testutil.IsSymlink(tt, fs, err == nil && !test.dry, test.to)
		})
	}
}

func TestUnlink(t *testing.T) {
	var tests = []struct {
		desc string
		file string
		dry  bool
		err  error
	}{
		{
			"Simple unlinking works",
			"a",
			false,
			nil,
		},
		{
			"It doesn't unlink when doing a dry run",
			"a",
			true,
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			fs := tempfs.New(osfs.New())
			defer fs.Cleanup()

			_, err := fs.Create(test.file)
			assert.Equal(tt, test.err, err)

			err = fileutil.Unlink(fs, test.file, test.dry)
			assert.Equal(tt, test.err, err)
			testutil.FileExists(tt, fs, err == nil && test.dry, test.file)
		})
	}
}
