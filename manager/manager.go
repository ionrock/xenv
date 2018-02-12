package manager

import (
	"errors"
	"fmt"
	"sync"
	"syscall"

	"github.com/codeskyblue/kexec"
	"github.com/ionrock/xenv/util"
)

// Manager manages a set of Processes.
type Manager struct {
	Processes map[string]*kexec.KCommand

	pipeWaits map[string]*sync.WaitGroup
	lock      sync.Mutex
}

// New creates a new *Manager.
func New() *Manager {
	return &Manager{
		Processes: make(map[string]*kexec.KCommand),
		pipeWaits: make(map[string]*sync.WaitGroup),
	}

}

// StdoutHandler returns an OutHandler that will ensure the underlying
// process has an empty stdout buffer and logs to stdout a prefixed value
// of "$name | $line".
func (m *Manager) StdoutHandler(name string) util.OutHandler {
	return func(line string) string {
		fmt.Printf("%s | %s\n", name, line)
		return ""
	}
}

// StderrHandler returns an OutHandler that will ensure the underlying
// process has an empty stderr buffer and logs to stdout a prefixed value
// of "$name | $line".
func (m *Manager) StderrHandler(name string) util.OutHandler {
	return func(line string) string {
		fmt.Printf("%s | %s\n", name, line)
		return ""
	}
}

// Start and managed a new process using the default handlers from a
// string.
func (m *Manager) Start(name, command, dir string, env []string) error {
	cmd := kexec.Command("/bin/bash", "-c", command)
	cmd.Dir = dir
	cmd.Env = env

	return m.StartProcess(name, cmd)
}

// StartProcess starts and manages a predifined process.
func (m *Manager) StartProcess(name string, cmd *kexec.KCommand) error {
	return m.StartProcessWithHandlers(name, cmd, m.StdoutHandler(name), m.StderrHandler(name))
}

// StartProcessWithHandlers
func (m *Manager) StartProcessWithHandlers(name string, cmd *kexec.KCommand, o, e util.OutHandler) error {
	m.lock.Lock()
	defer m.lock.Unlock()

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
	go util.LineReader(wg, stdout, o)
	go util.LineReader(wg, stderr, e)

	err = cmd.Start()
	if err != nil {
		return err
	}

	m.Processes[name] = cmd
	m.pipeWaits[name] = wg

	return nil

}

// Stop will try to stop a managed process. If the process does not
// exist, no error is returned.
func (m *Manager) Stop(name string) error {
	cmd, ok := m.Processes[name]
	// We don't mind stopping a process that doesn't exist.
	if !ok {
		return nil
	}

	// ProcessState means it is already exited.
	if cmd.ProcessState != nil {
		return nil
	}

	err := cmd.Terminate(syscall.SIGKILL)
	if err != nil {
		fmt.Println("Unable to kill process")
		return err
	}

	return nil
}

// Remove will try to stop and remove a managed process.
func (m *Manager) Remove(name string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	err := m.Stop(name)
	if err != nil {
		return err
	}

	// Note that if the stop fails we don't remove it from the map of
	// processes to avoid losing the reference.
	delete(m.Processes, name)

	return nil
}

// Wait will block until all managed processes have finished.
func (m *Manager) Wait() error {
	wg := &sync.WaitGroup{}
	wg.Add(len(m.Processes))

	for _, cmd := range m.Processes {
		go func(cmd *kexec.KCommand) {
			defer wg.Done()
			cmd.Wait()
		}(cmd)
	}

	wg.Wait()

	return nil
}

// WaifFor will wait on a specific process by name
func (m *Manager) WaitFor(name string) error {
	cmd, ok := m.Processes[name]
	if !ok {
		return errors.New("missing process")
	}

	return cmd.Wait()
}
