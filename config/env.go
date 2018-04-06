// Package config takes the xenv config and performs the tasks
// defined.
package config

import (
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codeskyblue/kexec"
	"github.com/ionrock/xenv/manager"
	"github.com/ionrock/xenv/util"
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

// Environment maintains the executable environment state.
type Environment struct {
	// Services provides a simple process manager to start/stop
	// processes along with the primary command.
	Services *manager.Manager

	// Tasks can run before and after a command.
	Tasks map[string]*exec.Cmd

	// Config provides the environment for the command.
	Config *Config

	// ConfigDir is the directory where the config file is in order to
	// provide a base for tasks / services.
	ConfigDir  string
	ConfigFile string

	DataOnly bool
	post     []*XeConfig
}

// NewEnvironment creates a new *Environment rooted at the provided
// directory.
func NewEnvironment() *Environment {
	return &Environment{
		Services: manager.New(),
		Tasks:    make(map[string]*exec.Cmd),
		Config:   &Config{make(map[string]string)},
	}
}

// NewEnvironmentFromConfig
func NewEnvironmentFromConfig(cfgFile string) (*Environment, error) {
	cfgDir, err := util.AbsDir(cfgFile)
	if err != nil {
		return nil, err
	}

	env := NewEnvironment()
	env.ConfigDir = cfgDir
	env.ConfigFile = cfgFile
	return env, nil
}

// Pre runs the defined steps before the specified command.
func (e *Environment) Pre() error {
	cfgs, err := e.Load()
	if err != nil {
		return err
	}

	for _, cfg := range cfgs {
		if err := e.ConfigHandler(cfg); err != nil {
			log.WithError(err).WithFields(log.Fields{
				"config": cfg,
			}).Warn("error running config")
			return err
		}
	}

	return nil
}

// Post runs the defined steps after the process exits, no matter the
// exit status of the command.
func (e *Environment) Post() error {
	if len(e.post) == 0 {
		return nil
	}

	log.Info("Running Post Events")
	for _, cfg := range e.post {
		// We don't worry about using a data handler here.
		if err := e.ConfigHandler(cfg); err != nil {
			return err
		}
	}

	return nil
}

// SetEnvFromEnvvars sets environment values from a list of key value
// pairs ([]map[string]string).
func (e *Environment) SetEnvFromEnvvars(envvars []map[string]string) error {
	for _, m := range envvars {
		for k, v := range m {
			err := e.SetEnv(k, v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// SetEnv sets an environment value.
func (e *Environment) SetEnv(k, v string) error {
	v = os.Expand(v, e.Config.GetConfig)
	val, err := CompileValue(v, e.ConfigDir, e.Config.ToEnv())

	if err != nil {
		log.WithFields(log.Fields{
			"value": v,
		}).WithError(err).Warn("error getting value for env")
		return err
	}

	log.WithFields(log.Fields{
		"key": k, "value": v,
	}).Debug("setting value")

	e.Config.Set(k, val)
	return nil
}

// SetEnvFromScript will run a script that outputs YAML or JSON,
// flatten the output and add it to the environment's configuration.
func (e *Environment) SetEnvFromScript(cmd, dir string) error {
	s := Script{
		Cmd: cmd,
		Dir: dir,
		Env: e.Config.ToEnv(),
	}

	env, err := s.Load()
	if err != nil {
		return err
	}

	for k, v := range env {
		// We expand the value if it has any vars defined. This will
		// also remove expansions that don't exist leaving things with an
		// empty string.
		val := os.Expand(v, e.Config.GetConfig)
		e.SetEnv(k, val)
	}

	return nil
}

// RunTask runs a task in the environment. The output is sent to
// stdout and is prefixed by the name of the task.
func (e *Environment) RunTask(name, command, dir string) error {
	if dir == "" {
		dir = e.ConfigDir
	}

	t := &Task{
		Name: name,
		Cmd:  command,
		Dir:  dir,
		Env:  e.Config.ToEnv(),
	}

	return t.Run()
}

// ConfigHandler calls the respective handler actionss based on the
// passed in XeConfig. It is assumed the XeConfig will only have 1 field
// in its struct filled in.
func (e *Environment) ConfigHandler(cfg *XeConfig) error {
	switch {
	case cfg.Env != nil:
		err := e.SetEnvFromEnvvars(cfg.Env)
		if err != nil {
			return err
		}

	case cfg.EnvScript != "":
		err := e.SetEnvFromScript(cfg.EnvScript, e.ConfigDir)
		if err != nil {
			return err
		}

	case cfg.Template != nil && !e.DataOnly:
		cfg.Template.Env = e.Config.Data
		err := cfg.Template.Execute(e.ConfigDir)
		if err != nil {
			return err
		}

	case cfg.Task != nil && !e.DataOnly:
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

// StopServices stops the services managed by the process manager.
func (e *Environment) StopServices() error {
	for name := range e.Services.Processes {
		err := e.Services.Stop(name)
		if err != nil {
			log.WithFields(log.Fields{
				"service_name": name,
				"error":        err,
			}).Error("problem stopping service")
			return err
		}
	}

	return nil
}

func (e *Environment) Load() ([]*XeConfig, error) {
	log.Debugf("loading %s", e.ConfigFile)
	cfgs, err := NewXeConfig(e.ConfigFile)
	if err != nil {
		log.WithFields(log.Fields{
			"config_file": e.ConfigFile,
			"config_dir":  e.ConfigDir,
		}).WithError(err).Error("error loading config")
		return nil, err
	}
	return cfgs, nil
}

func (e *Environment) watchSignals(done chan error, cmd *kexec.KCommand) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		err := cmd.Terminate(sig)
		if err != nil {
			log.WithError(err).Warn("error sending signal")
		}
		done <- err
	}()
}

func (e *Environment) watchConfig(cadence *time.Ticker, events chan struct{}) {
	go func() {
		for range cadence.C {
			ne, err := NewEnvironmentFromConfig(e.ConfigFile)
			if err != nil {
				log.WithError(err).Warn("error rebuilding config data")
				continue
			}

			ne.DataOnly = true
			err = ne.Pre()
			if err != nil {
				log.WithError(err).Warn("error rebuilding config data")
				continue
			}

			// compare the envs
			diff := e.Config.Diff(ne.Config)
			if diff != nil {
				fields := log.Fields{}
				for k, v := range diff.Data {
					fields[k] = v
				}
				log.WithFields(fields).Debug("env diff")
				events <- restartEvent{}
			}
		}
	}()

}

func (e *Environment) wait(cmd *kexec.KCommand, done chan error, events chan struct{}) error {
	for {
		select {
		case err := <-done:
			if err != nil {
				log.WithError(err).Warn("process exit error")
			}
			return err
		case <-events:
			log.Info("Configuration data changed. Exiting...")

			err := cmd.Terminate(syscall.SIGINT)
			if err != nil {
				log.WithError(err).Warn("error stopping process")
				if cmd.ProcessState == nil {
					err = cmd.Terminate(syscall.SIGTERM)
					if err != nil {
						log.WithError(err).Warn("error killing process")
					}
				}
			}
			return err
		}
	}

	return nil
}

// Main runs the configuration items, the main process and any post processes.
func (e *Environment) Main(parts []string) (err error) {
	err = e.Pre()
	if err != nil {
		return err
	}

	if len(parts) == 0 {
		return nil
	}

	// replace any replacements
	for i := range parts {
		parts[i] = os.Expand(parts[i], e.Config.GetConfig)
	}

	log.Infof("Running command: %s", strings.Join(parts, " "))

	cmd := kexec.Command(parts[0])
	if len(parts) > 1 {
		cmd.Args = append(cmd.Args, parts[1:]...)
	}
	cmd.Env = e.Config.ToEnv()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start our process and listen for signals
	done := make(chan error)
	e.watchSignals(done, cmd)

	go func() {
		done <- cmd.Run()
	}()

	events := make(chan struct{})
	cadence := time.NewTicker(10 * time.Second)
	defer cadence.Stop()

	e.watchConfig(cadence, events)

	err = e.wait(cmd, done, events)

	postErr := e.Post()
	if postErr != nil {
		log.WithError(postErr).Warn("Error running post")
	}

	return err
}

type restartEvent struct{}
