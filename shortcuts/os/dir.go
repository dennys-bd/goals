package os

import (
	"io"
	"os"

	errs "github.com/dennys-bd/goals/shortcuts/errors"
)

// Exists verify is a file or dir exists
func Exists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if !os.IsNotExist(err) {
		errs.Ex(err)
	}
	return false
}

// IsEmpty checks if a given path is empty.
// Hidden files in path are ignored.
func IsEmpty(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		errs.Ex(err)
	}

	if !fi.IsDir() {
		return fi.Size() == 0
	}

	f, err := os.Open(path)
	if err != nil {
		errs.Ex(err)
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil && err != io.EOF {
		errs.Ex(err)
	}

	for _, name := range names {
		if len(name) > 0 && name[0] != '.' {
			return false
		}
	}
	return true
}
