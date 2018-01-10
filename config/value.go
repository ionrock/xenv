package config

import (
	"bytes"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/ionrock/xenv/util"
)

// CompileValue accepts a value and tries to run it as a command when
// it starts and ends with backticks (`). The path provides the directory
// where the command will be run and the []string, the environment. The
// result will be trimmed of whitespace in order to be used as a string
// value. For example, if a command normally would output an extra new
// line for the terminal, that newline is removed.
func CompileValue(value, path string, env []string) (string, error) {
	log.Debug("%#v", value)

	if !strings.HasPrefix(value, "`") || !strings.HasSuffix(value, "`") {
		return value, nil
	}

	dirname, err := util.AbsDir(path)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("/bin/bash", "-c", strings.Trim(value, "`"))
	cmd.Dir = dirname
	if len(env) > 0 {
		cmd.Env = env
	}

	buf, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(buf)), nil
}
