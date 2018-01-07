package config

import (
	"bytes"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/ionrock/xenv/util"
)

func CompileValue(value, path string, env []string) (string, error) {
	log.Debug("%#v", value)

	if !strings.HasPrefix(value, "`") || !strings.HasSuffix(value, "`") {
		return value, nil
	}

	dirname, err := util.AbsDir(path)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("sh", "-c", strings.Trim(value, "`"))
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
