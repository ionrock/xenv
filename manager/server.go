package manager

import (
	fmt "fmt"

	context "golang.org/x/net/context"
)

type ManagerServer struct {
	Mgr *Manager
}

// Start should start the specified process. If the process is already started it is a noop.
func (srv *ManagerServer) Start(ctx context.Context, proc *SvcProcess) (*SvcResult, error) {
	return &SvcResult{fmt.Sprintf("started %s", proc.Name)}, nil
}

// Stop this should stop the process but the manager needs some code for that...
func (srv *ManagerServer) Stop(ctx context.Context, proc *SvcProcess) (*SvcResult, error) {
	err := srv.Mgr.Stop(proc.Name)
	if err != nil {
		return nil, err
	}
	return &SvcResult{fmt.Sprintf("stopped %s", proc.Name)}, nil
}

// Restart requests the manager to restart the process.
func (srv *ManagerServer) Restart(ctx context.Context, proc *SvcProcess) (*SvcResult, error) {
	err := srv.Mgr.Stop(proc.Name)
	if err != nil {
		return nil, err
	}
	return &SvcResult{fmt.Sprintf("restarted %s", proc.Name)}, nil
}
