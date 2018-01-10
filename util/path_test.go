package util_test

import (
	"path/filepath"
	"testing"

	"github.com/ionrock/xenv/util"
)

func TestAbsDir(t *testing.T) {
	a, err := util.AbsDir("testdata/config.yml")
	if err != nil {
		t.Errorf("unexpected error with absdir: %s", err)
	}

	abs, err := filepath.Abs(a)
	if err != nil {
		t.Errorf("unexpected error with abspath: %s", err)
	}

	if abs != a {
		t.Errorf("invalid absdir: %s != %s", abs, a)
	}
}
