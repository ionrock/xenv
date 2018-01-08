package process

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/ionrock/xenv/util"
)

type Process struct {
	Cmd       *exec.Cmd
	Output    *Output
	pipesWait *sync.WaitGroup
}

func New(command string, dir string) *Process {
	var cmd *exec.Cmd
	parts := util.SplitCommand(command)
	if len(parts) == 0 {
		cmd = exec.Command(parts[0])
	} else {
		cmd = exec.Command(parts[0], parts[1:]...)
	}

	cmd.Dir = dir

	return &Process{
		Cmd:       cmd,
		pipesWait: new(sync.WaitGroup),
	}
}

// Run the exec.Cmd handling stdout / stderr according to the
// configured Output struct.
func (p *Process) Run() error {
	if p.Output == nil {
		return p.Cmd.Run()
	}

	err := p.Start()
	if err != nil {
		p.Output.SystemOutput(fmt.Sprint("Failed to start ", p.Cmd.Args, ": ", err))
		return err
	}

	return p.Wait()
}

// Start the exec.Cmd start.
func (p *Process) Start() error {
	stdout, err := p.Cmd.StdoutPipe()
	if err != nil {
		fmt.Println("error creating stdout pipe")
		return err
	}
	stderr, err := p.Cmd.StderrPipe()
	if err != nil {
		fmt.Println("error creating stderr pipe")
		return err
	}

	if p.pipesWait == nil {
		p.pipesWait = new(sync.WaitGroup)
	}
	p.pipesWait.Add(2)

	go p.Output.LineReader(p.pipesWait, p.Output.Name, stdout, false)
	go p.Output.LineReader(p.pipesWait, p.Output.Name, stderr, true)

	return p.Cmd.Start()
}

func (p *Process) Wait() error {
	if p.pipesWait != nil {
		p.pipesWait.Wait()
	}
	return p.Cmd.Wait()
}

func (p *Process) Terminate() error {
	if p.Cmd.ProcessState != nil {
		return nil
	}

	return p.Cmd.Process.Kill()
}

func ParseEnv(environ []string) map[string]string {
	env := make(map[string]string)
	for _, e := range environ {
		pair := strings.SplitN(e, "=", 2)
		env[pair[0]] = pair[1]
	}
	return env
}

func Env(env map[string]string, useEnv bool) []string {
	envlist := []string{}

	// update our env by loading our env and overriding any values in
	// the provided env.
	if useEnv {
		environ := ParseEnv(os.Environ())
		for k, v := range env {
			environ[k] = v
		}
		env = environ
	}

	for key, val := range env {
		if key == "" {
			continue
		}
		envlist = append(envlist, fmt.Sprintf("%s=%s", key, val))
	}

	return envlist
}