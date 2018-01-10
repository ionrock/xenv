package util

import (
	"os"
	"path/filepath"
)

// AbsDir finds the absolute directory name of a path (string) and
// asserts the file exists.
func AbsDir(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	f, err := os.Stat(abs)
	if err != nil {
		return "", err
	}

	if f.IsDir() {
		return abs, nil
	}

	return filepath.Dir(abs), nil
}
