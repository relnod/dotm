package file

import (
	"log"
	"os"
)

// IsSymlink checks if the given path is a symlink.
func IsSymlink(path string) bool {
	_, err := os.Lstat(path)
	if err != nil {
		return false
	}
	_, err = os.Readlink(path)
	if err != nil {
		return false
	}
	return true
}

// Link creates a symbolik link from source to dest. When dry is true only
// perfomers a dry run, by printing the perfomed action.
func Link(source string, dest string, dry bool) error {
	if dry {
		log.Printf("Creating symlink: %s -> %s\n", source, dest)
		return nil
	}

	err := os.Symlink(source, dest)
	if err != nil {
		return err
	}

	return nil
}

// Unlink removes a symbolik link from the given filepath. When dry is true only
// perfomers a dry run, by printing the performed action.
func Unlink(filepath string, dry bool) error {
	if dry {
		log.Printf("Removing symlink: %s\n", filepath)
		return nil
	}
	err := os.Remove(filepath)
	if err != nil {
		return err
	}

	return nil
}
