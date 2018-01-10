package config

import (
	"fmt"
	"os/exec"
	"sync"

	"github.com/ionrock/xenv/util"
)

// Task	runs a command before starting the process.
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

	fmt.Println("Running Task: " + t.Name)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("error creating stdout pipe: %s\n", err)
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("error creating stderr pipe: %s\n", err)
		return err
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)

	// These close the stdout/err channels
	if t.StdoutHandler == nil {
		t.StdoutHandler = t.outhandler
	}

	if t.StderrHandler == nil {
		t.StdoutHandler = t.outhandler
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

func (t *Task) outhandler(line string) string {
	fmt.Printf("%s | %s\n", t.Name, line)
	return line
}
