package util

import (
	"os"
	"path/filepath"
)

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
