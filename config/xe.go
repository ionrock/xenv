package config

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/ionrock/xenv/templates"
)

// Service is a service xenv config.
type Service struct {
	Name string `json:"name"`
	Cmd  string `json:"cmd"`
	Dir  string `json:"dir"`
}

// XeTask is a task in a xenv config.
type XeTask struct {
	Name string `json:"name"`
	Cmd  string `json:"cmd"`
	Dir  string `json:"dir"`
}

// the post in a xenv config.
type XeConfig struct {
	Service   *Service            `json:"service"`
	Env       []map[string]string `json:"env"`
	EnvScript string              `json:"envscript"`
	Task      *XeTask             `json:"task"`
	Post      []*XeConfig         `json:"post"`
	Template  *templates.Renderer `json:"template"`
}

// NewXeConfig parses a path for a *XeConfig.
func NewXeConfig(path string) ([]*XeConfig, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := make([]*XeConfig, 0)

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
