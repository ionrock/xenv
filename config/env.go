package config

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ionrock/xenv/manager"
)

func findLongestServiceName(cfgs []*XeConfig) int {
	size := 0

	for _, cfg := range cfgs {
		if cfg.Service == nil {
			continue
		}

		if len(cfg.Service.Name) > size {
			size = len(cfg.Service.Name)
		}
	}

	return size
}

type Environment struct {
	Services  *manager.Manager
	Tasks     map[string]*exec.Cmd
	Config    *Config
	ConfigDir string
	DataOnly  bool

	post []*XeConfig
}

func NewEnvironment(cfgDir string) *Environment {
	return &Environment{
		Services:  manager.New(),
		Tasks:     make(map[string]*exec.Cmd),
		Config:    &Config{make(map[string]string)},
		ConfigDir: cfgDir,
	}
}

func (e *Environment) Pre(cfgs []*XeConfig) error {
	handler := e.ConfigHandler
	if e.DataOnly {
		handler = e.DataHandler
	}

	for _, cfg := range cfgs {
		if err := handler(cfg); err != nil {
			return err
		}
	}

	return nil
}

func (e *Environment) Post() error {
	fmt.Println(fmt.Sprintf("post value: %#v", e.post))
	err := e.CleanUp()
	if err != nil {
		return err
	}

	for _, cfg := range e.post {
		// We don't worry about using a data handler here.
		if err := e.ConfigHandler(cfg); err != nil {
			return err
		}
	}

	return nil
}

func (e *Environment) StartService(name, command, dir string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = dir
	cmd.Env = e.Config.ToEnv()

	return e.Services.StartProcess(name, cmd)
}

func (e *Environment) SetEnv(k, v string) error {
	v = os.Expand(v, e.Config.GetConfig)
	val, err := CompileValue(v, e.ConfigDir, e.Config.ToEnv())
	if err != nil {
		fmt.Printf("error getting value for env: %q %q\n", v, err)
		return err
	}
	e.Config.Set(k, val)
	return nil
}

func (e *Environment) SetEnvFromScript(cmd, dir string) error {
	s := Script{Cmd: cmd, Dir: dir}
	err := s.Apply(e.Config)
	if err != nil {
		return err
	}

	return nil
}

func (e *Environment) RunTask(name, command, dir string) error {
	t := &Task{
		Name: name,
		Cmd:  command,
		Dir:  dir,
		Env:  e.Config.ToEnv(),
	}

	return t.Run()
}

func (e *Environment) DataHandler(cfg *XeConfig) error {
	switch {
	case cfg.Env != nil:
		for k, v := range cfg.Env {
			err := e.SetEnv(k, v)
			return err
		}

	case cfg.EnvScript != "":
		err := e.SetEnvFromScript(cfg.EnvScript, e.ConfigDir)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Environment) ConfigHandler(cfg *XeConfig) error {
	switch {
	case cfg.Service != nil:
		if cfg.Service.Dir == "" {
			cfg.Service.Dir = e.ConfigDir
		}

		err := e.StartService(cfg.Service.Name, cfg.Service.Cmd, cfg.Service.Dir)
		if err != nil {
			return err
		}

	case cfg.Env != nil:
		for k, v := range cfg.Env {
			err := e.SetEnv(k, v)
			return err
		}

	case cfg.EnvScript != "":
		err := e.SetEnvFromScript(cfg.EnvScript, e.ConfigDir)
		if err != nil {
			return err
		}

	case cfg.Task != nil:
		err := e.RunTask(cfg.Task.Name, cfg.Task.Cmd, cfg.Task.Dir)
		if err != nil {
			return err
		}

	case cfg.Post != nil:
		if e.post == nil {
			e.post = make([]*XeConfig, 0)
		}

		for _, cfg := range cfg.Post {
			e.post = append(e.post, cfg)
		}
	}

	return nil
}

func (e *Environment) CleanUp() error {
	for name, _ := range e.Services.Processes {
		err := e.Services.Stop(name)
		if err != nil {
			return err
		}
	}

	return nil
}
