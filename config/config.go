package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Data map[string]string
}

func (c *Config) GetConfig(name string) string {
	if v, ok := c.Data[name]; ok {
		return v
	}

	return os.Getenv(name)
}

func (c *Config) Set(k, v string) {
	c.Data[k] = v
}

func (c *Config) Get(k string) (string, bool) {
	v, ok := c.Data[k]
	return v, ok
}

// ToEnv returns the config data as a list of strings that can be used
// in exec.Cmd
func (c *Config) ToEnv() []string {
	envlist := []string{}
	for key, val := range c.Data {
		if key == "" {
			continue
		}

		if val == "" && os.Getenv(key) != "" {
			val = os.Getenv(key)
		}

		envlist = append(envlist, fmt.Sprintf("%s=%s", key, val))
	}

	return envlist
}

var osEnviron = os.Environ

// Environ overlays the config data on top of the existing process
// environment.
func (c *Config) Environ() []string {
	envlist := []string{}

	// filter out values that we have in our config data
	for _, envvar := range osEnviron() {
		key := strings.SplitN(envvar, "=", 2)[0]
		if _, ok := c.Data[key]; !ok {
			envlist = append(envlist, envvar)
		}
	}

	// Then add our config
	envlist = append(envlist, c.ToEnv()...)

	return envlist
}
