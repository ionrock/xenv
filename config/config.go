package config

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Config provides a managed map that provides configuration for an Environment.
type Config struct {
	Data map[string]string
}

// GetConfig gets a config value from the Config, falling back to the
// os.Environ. This function can be used with os.Expand.
func (c *Config) GetConfig(name string) string {
	if v, ok := c.Data[name]; ok {
		return v
	}

	return os.Getenv(name)
}

// Set sets a value in the Config
func (c *Config) Set(k, v string) {
	c.Data[k] = v
}

// Get gets a value in the config and is compatible with a map.
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

	sort.Strings(envlist)
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

func (c *Config) Diff(o *Config) *Config {
	diff := &Config{make(map[string]string)}

	compareConfigs(c, o, diff)
	compareConfigs(o, c, diff)

	if len(diff.Data) > 0 {
		return diff
	}

	return nil
}

func compareConfigs(a, b, diff *Config) {
	for key, val := range a.Data {
		otherVal, ok := b.Get(key)
		if !ok {
			diff.Set(key, val)
		}
		if otherVal != val {
			diff.Set(key, otherVal)
		}
	}
}
