package test

import (
	"testing"

	"github.com/relnod/fsa"
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
			fs := fsa.NewTempFs(fsa.NewOsFs())
			defer fs.Cleanup()

			assert.NoError(tt, testutil.CreateFiles(fs, c.given))
			assert.NoError(tt, testutil.AddFiles(fs, "./testdata", "/"))
			assert.NoError(tt, testutil.PrintFiles(fs))
			assert.NoError(tt, c.exec(fs.Base()))
			assert.NoError(tt, testutil.CheckFiles(fs, c.expected))
		})
	}
}
