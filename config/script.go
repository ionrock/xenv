package config

import (
	"os/exec"

	"github.com/ghodss/yaml"
)

// Script runs a script and tries to parse it as YAML or JSON for
// application to the Environment.
type Script struct {
	Cmd string
	Dir string
	Env []string
}

// Load executes the script using the specified *Config for the
// environment. The result is parsed and flattened before returning a
// map[string]string.
func (e Script) Load() (map[string]string, error) {
	cmd := exec.Command("/bin/bash", "-c", e.Cmd)
	cmd.Dir = e.Dir
	cmd.Env = e.Env

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
