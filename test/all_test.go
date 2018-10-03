package test

import (
	"testing"

	"github.com/relnod/fsa/osfs"
	"github.com/relnod/fsa/tempfs"
	"github.com/relnod/fsa/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	testcases, err := parseDir("./testcases")
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range testcases {
		t.Run(c.name, func(tt *testing.T) {
			fs := tempfs.New(osfs.New())
			defer fs.Cleanup()

			assert.NoError(tt, testutil.CreateFiles(fs, c.given))

			assert.NoError(tt, c.exec(fs.Base()))

			assert.NoError(tt, testutil.CheckFiles(fs, c.expected))
		})
	}
}
