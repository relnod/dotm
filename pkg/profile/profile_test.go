package profile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandPath(t *testing.T) {
	home := os.Getenv("HOME")
	defer os.Setenv("HOME", home)

	os.Setenv("HOME", "/testhome")
	currDir, _ := filepath.Abs(".")

	tests := []struct {
		given string
		want  string
	}{
		{"$HOME/.dotfiles/<PROFILE>/", "/testhome/.dotfiles/foo/"},
		{".", currDir},
	}

	for _, test := range tests {
		t.Run(test.given, func(tt *testing.T) {
			got, err := expandPath(test.given, "foo")
			assert.NoError(tt, err)
			assert.Equal(tt, test.want, got)
		})
	}
}
