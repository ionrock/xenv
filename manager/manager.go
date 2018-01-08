package manager

import (
	"fmt"
	"os/exec"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/ionrock/xenv/process"
	"github.com/ionrock/xenv/util"
)

type Manager struct {
	// processes is a map with our processes
	ProcessMap map[string]*exec.Cmd

	// Output provides a prefix formatter for logging
	Output *process.Output

	// wg for keeping track of our process go routines
	wg sync.WaitGroup

	teardown, teardownNow Barrier
}

func New(of *process.Output) *Manager {
	return &Manager{
		ProcessMap: make(map[string]*exec.Cmd),
		Output:     of,
	}
}

func (m *Manager) Processes() map[string]*exec.Cmd {
	return m.ProcessMap
}

func (m *Manager) Start(name, command, dir string, env []string, of *process.Output) error {
	if of == nil {
		of = m.Output
	}

	parts := util.SplitCommand(command)

	var ps *exec.Cmd
	if len(parts) == 1 {
		ps = exec.Command(parts[0])
	} else {
		ps = exec.Command(parts[0], parts[1:]...)
	}

	if dir != "" {
		ps.Dir = dir
	}

	if env != nil {
		ps.Env = env
	}

	stdout, err := ps.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := ps.StderrPipe()
	if err != nil {
		return err
	}

	// Start reading the output of the
	pipeWait := new(sync.WaitGroup)
	pipeWait.Add(2)

	go of.LineReader(pipeWait, name, stdout, false)
	go of.LineReader(pipeWait, name, stderr, true)

	finished := make(chan struct{}) // closed on process exit

	err = ps.Start()
	if err != nil {
		of.SystemOutput(fmt.Sprint("Failed to start ", name, ": ", err))
		return err
	}

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		defer close(finished)
		pipeWait.Wait()
		ps.Wait()
	}()

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()

		select {
		case <-finished:
			m.teardown.Fall()

		case <-m.teardown.Barrier():
			of.SystemOutput(fmt.Sprintf("Killing %s", name))
			ps.Process.Kill()
		}
	}()

	m.ProcessMap[name] = ps

	return nil
}

func (m *Manager) Stop(name string) error {
	svc, ok := m.ProcessMap[name]
	if !ok {
		// should probably still throw an error here...
		return nil
	}

	if svc.ProcessState != nil {
		return nil
	}

	err := svc.Process.Kill()
	if err != nil {
		log.Printf("error killing service: %s, %s\n", name, err)
		return err
	}

	return nil
}
