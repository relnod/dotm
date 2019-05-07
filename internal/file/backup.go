package file

import (
	"io/ioutil"
	"log"
	"os"
)

const backupFileSuffix = ".backup"

func backupPath(path string) string {
	return path + backupFileSuffix
}

// Backup moves the given file to "filename.backup". When dry is true only
// perfomers a dry run, by printing the performed action.
func Backup(file string, dry bool) error {
	if dry {
		log.Printf("Creating backup: %s\n", file)
		return nil
	}
	return moveFile(file, backupPath(file))
}

// RestoreBackup tries to restore a backup. When dry is true only
// perfomers a dry run, by printing the performed action.
func RestoreBackup(file string, dry bool) error {
	if _, err := os.Stat(backupPath(file)); os.IsNotExist(err) {
		return nil
	}
	if dry {
		log.Printf("Restoring backup: %s\n", file)
		return nil
	}
	return moveFile(backupPath(file), file)
}

func moveFile(source, dest string) error {
	data, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dest, data, os.ModePerm)
	if err != nil {
		return err
	}
	return os.Remove(source)
}
