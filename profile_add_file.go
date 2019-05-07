package dotm

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/relnod/dotm/internal/file"
)

func (p *Profile) addFile(dir, path string) error {
	sourceFile := filepath.Join(os.Getenv("HOME"), path)
	destFile := filepath.Join(p.expandedPath, dir, path)
	err := os.MkdirAll(filepath.Dir(destFile), os.ModePerm)
	if err != nil {
		return err
	}
	data := []byte("# Created by dotm")

	// Check if the source file already exists.
	if _, err := os.Stat(sourceFile); err == nil {
		if file.IsSymlink(sourceFile) {
			err = file.Unlink(sourceFile, false)
			if err != nil {
				return err
			}
		} else {
			data, err = ioutil.ReadFile(sourceFile)
			if err != nil {
				return err
			}
			err = file.Backup(sourceFile, false)
			if err != nil {
				return err
			}
		}
	}

	// Create the file, since it does not exist.
	err = ioutil.WriteFile(destFile, data, os.ModePerm)
	if err != nil {
		return err
	}

	// We link back from the destination file to the source file.
	return file.Link(destFile, sourceFile, false)
}
