package config

import (
	"os"
	"os/exec"

	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
	"github.com/ionrock/we/flat"
)

type Script struct {
	Cmd string
	Dir string
}

func (e Script) Load() (map[string]string, error) {
	cmd := exec.Command("sh", "-c", e.Cmd)
	cmd.Dir = e.Dir

	buf, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var f interface{}

	err = yaml.Unmarshal(buf, &f)
	if err != nil {
		return nil, err
	}

	env := &flat.FlatEnv{
		Path: e.Dir,
		Env:  make(map[string]string),
	}

	err = env.Load(f, []string{})
	if err != nil {
		return nil, err
	}

	return env.Env, nil
}

func (e Script) Apply(config *Config) error {
	env, err := e.Load()
	if err != nil {
		return err
	}

	for k, v := range env {
		log.Debugf("Setting: %s to %s", k, os.Expand(v, config.GetConfig))
		val, err := CompileValue(os.Expand(v, config.GetConfig), e.Dir)
		if err != nil {
			return err
		}
		config.Set(k, val)
	}

	return nil
}
