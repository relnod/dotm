package fsa

import (
	"math/rand"
	"path/filepath"
	"time"
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func NewTempFs(fs FileSystem) *BaseFs {
	return NewBaseFs(fs, randomTempDir(fs))
}

func randomTempDir(fs FileSystem) string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = "abcdefgh"[seededRand.Intn(8)]
	}
	path := filepath.Join(fs.TempDir(), string(b))
	if _, err := fs.Stat(path); err == nil {
		return randomTempDir(fs)
	}
	return path
}
