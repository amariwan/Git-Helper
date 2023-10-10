package helper

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func FileOrDirExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func DirExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	if stat.IsDir() {
		return true
	}

	return false
}

func FindUpwards(path string, fileName string) string {
	var basePath string = "/"
	for {
		rel, _ := filepath.Rel(basePath, path)

		// Exit the loop once we reach the basePath.
		if rel == "." {
			break
		}

		if FileOrDirExists(fmt.Sprintf("%v/%v", path, fileName)) {
			return path
		}

		// Going up!
		path += "/.."
	}

	return ""
}

type FindDownwardsFunc func(path string)

type walker struct {
	FileName string
	Fn       FindDownwardsFunc
}

func (w *walker) walk(s string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if !d.IsDir() && strings.HasSuffix(s, w.FileName) {
		dir, _ := path.Split(s)
		if w.Fn != nil {
			w.Fn(dir)
		}
	}
	return nil
}

func FindDownwards(startDir string, fileName string, fn FindDownwardsFunc) {
	var w = walker{FileName: fileName, Fn: fn}
	filepath.WalkDir(startDir, w.walk)
}
