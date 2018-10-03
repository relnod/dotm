package testutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/relnod/fsa"
	"github.com/relnod/fsa/fsutil"
)

type file struct {
	path      string
	isDir     bool
	isSymlink bool
	isDeleted bool
	content   string
}

func parse(raw string) []file {
	raw = os.ExpandEnv(raw)
	r := regexp.MustCompile("[^\\s]+")
	paths := r.FindAllString(raw, -1)
	var files []file
	for _, path := range paths {
		f := file{path: path}
		if strings.HasSuffix(path, "/") {
			f.isDir = true
		}
		modifier := strings.Split(path, ":")
		if len(modifier) == 2 {
			f.path = modifier[0]
			if modifier[1] == "ln" {
				f.isSymlink = true
			} else if modifier[1] == "deleted" {
				f.isDeleted = true
			}
		}
		content := strings.Split(path, "#")
		if len(content) == 2 {
			f.path = content[0]
			f.content = content[1]
		}

		files = append(files, f)
	}
	return files
}

func CreateFiles(fs fsa.FileSystem, raw string) error {
	for _, file := range parse(raw) {
		if file.isDir {
			err := fs.MkdirAll(file.path, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}
		err := createFile(fs, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func createFile(fs fsa.FileSystem, f file) error {
	if f.path == "" {
		return nil
	}
	s := strings.Split(f.path, "/")
	if len(s) > 1 {
		sb := strings.Builder{}
		for _, dir := range s[:len(s)-1] {
			sb.WriteString(dir)
			fs.MkdirAll(sb.String(), os.ModePerm)
			sb.WriteByte('/')
		}
	}

	if f.isSymlink {
		link := filepath.Join("/tmp", f.path)
		err := createFile(fs, file{path: link})
		if err != nil {
			return err
		}

		return fs.Symlink(link, f.path)
	}
	return fsutil.WriteFile(fs, f.path, []byte(f.content), os.ModePerm)
}

func CheckFiles(fs fsa.FileSystem, raw string) error {
	for _, file := range parse(raw) {
		if file.isDeleted {
			if FileExists(fs, file.path) {
				return fmt.Errorf("%s should be deleted", file.path)
			}
			continue
		}
		if !FileExists(fs, file.path) {
			return fmt.Errorf("%s doesn't exist", file.path)
		}
		if file.isDir {
			if !DirExists(fs, file.path) {
				return fmt.Errorf("%s should be a directory", file.path)
			}
			continue
		}
		if file.isSymlink && !IsSymlink(fs, file.path) {
			return fmt.Errorf("%s should be a symlink", file.path)
		}
	}
	return nil
}

// AddFiles adds all files from the given directory to the given file system.
// Works recusively.
func AddFiles(fs fsa.FileSystem, src, dest string) error {
	return addFiles(fs, src, dest)
}

func addFiles(fs fsa.FileSystem, src, dest string) error {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, f := range files {
		sourcePath := filepath.Join(src, f.Name())
		destPath := filepath.Join(dest, f.Name())
		if f.IsDir() {
			err = fs.MkdirAll(destPath, os.ModePerm)
			if err != nil {
				return err
			}

			err = addFiles(fs, sourcePath, destPath)
			if err != nil {
				return err
			}
			continue
		}

		data, err := ioutil.ReadFile(sourcePath)
		if err != nil {
			return err
		}

		err = fsutil.WriteFile(fs, destPath, data, os.ModePerm)
		if err != nil {
			return err
		}

	}

	return nil
}

// PrintFiles prints all files in the file system to stdout.
func PrintFiles(fs fsa.FileSystem) error {
	return printFiles(fs, "/")
}

func printFiles(fs fsa.FileSystem, dir string) error {
	files, err := fsutil.ReadDir(fs, dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		path := filepath.Join(dir, f.Name())
		if f.IsDir() {
			err := printFiles(fs, path)
			if err != nil {
				return err
			}
			continue
		}
		fmt.Println(path)
	}

	return nil
}
