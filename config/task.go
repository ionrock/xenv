package config

import (
	"os/exec"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/ionrock/xenv/util"
)

// Task runs a command before starting the process.
type Task struct {
	// Name the task. This Name is used to prefix the output to stdout.
	Name string

	// Cmd is the command to execute. This will be run in a sh.
	Cmd string

	// Dir is the directory to run the command.
	Dir string

	// Env is the environment to use for the command.
	Env []string

	StdoutHandler util.OutHandler
	StderrHandler util.OutHandler
}

// Run runs the command and prints the output to stdout prefixed by the Name.
func (t *Task) Run() error {
	cmd := exec.Command("/bin/bash", "-c", t.Cmd)
	cmd.Dir = t.Dir
	cmd.Env = t.Env

	name := t.Name
	if name == "" {
		name = t.Cmd
	}
	taskLog := log.WithFields(log.Fields{"name": name})

	taskLog.Info("Running Task")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.WithError(err).Printf("error creating stdout pipe")
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.WithError(err).Printf("error creating stderr pipe")
		return err
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)

	outhandler := func(line string) string {
		taskLog.Info(line)
		return line
	}

	// These close the stdout/err channels
	if t.StdoutHandler == nil {
		t.StdoutHandler = outhandler
	}

	if t.StderrHandler == nil {
		t.StderrHandler = outhandler
	}
	go util.LineReader(wg, stdout, t.StdoutHandler)
	go util.LineReader(wg, stderr, t.StderrHandler)

	err = cmd.Start()
	if err != nil {
		return err
	}

	wg.Wait()

	return cmd.Wait()
}
