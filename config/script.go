package config

import (
	"os"
	"os/exec"

	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
)

type Script struct {
	Cmd string
	Dir string
}

func (e Script) Load(config *Config) (map[string]string, error) {
	cmd := exec.Command("sh", "-c", e.Cmd)
	cmd.Dir = e.Dir
	cmd.Env = config.ToEnv()

	buf, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var f interface{}

	err = yaml.Unmarshal(buf, &f)
	if err != nil {
		return nil, err
	}

	env := &FlatEnv{
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
	env, err := e.Load(config)
	if err != nil {
		return err
	}

	for k, v := range env {
		log.Debugf("Setting: %s to %s", k, os.Expand(v, config.GetConfig))
		val, err := CompileValue(os.Expand(v, config.GetConfig), e.Dir, config.ToEnv())
		if err != nil {
			return err
		}
		config.Set(k, val)
	}

	return nil
}
