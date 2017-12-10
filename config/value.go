package config

import (
	"bytes"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/ionrock/xenv/process"
	"github.com/ionrock/xenv/util"
)

func CompileValue(value string, path string) (string, error) {
	log.Debug("%#v", value)

	if !strings.HasPrefix(value, "`") || !strings.HasSuffix(value, "`") {
		return value, nil
	}

	dirname, err := util.AbsDir(path)
	if err != nil {
		return "", err
	}

	proc := process.NewScript(strings.Trim(value, "`"), dirname)

	buf, err := proc.Execute()
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(buf.Bytes())), nil
}
