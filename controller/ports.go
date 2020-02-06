package controller

import (
	"errors"
	"fmt"
)

// AllocateAny associates a worker with any RDTP port
func (ctrl *Controller) AllocateAny(w *Worker) (uint16, error) {
	ctrl.Lock()
	defer ctrl.Unlock()

	for port := uint16(1); port < maxPortNo; port++ {
		if _, ok := ctrl.Ports[port]; !ok {
			ctrl.Ports[port] = w // reserve port for worker
			w.id = port          // assign port to worker
			return port, nil
		}
	}

	return 0, fmt.Errorf("all ports in use")
}

// Allocate associates a worker with a given RDTP port
func (ctrl *Controller) Allocate(w *Worker, p uint16) error {
	ctrl.Lock()
	defer ctrl.Unlock()

	if p == 0 {
		return errors.New("port 0 cannot be used (reserved)")
	}

	if _, ok := ctrl.Ports[p]; !ok {
		ctrl.Ports[p] = w // reserve port for worker
		w.id = p          // assign port to worker
		return nil
	}

	return errors.New("port is in use")
}

// Deallocate frees up an RDTP port
func (ctrl *Controller) Deallocate(port uint16) {
	ctrl.Lock()
	defer ctrl.Unlock()

	delete(ctrl.Ports, port)
}
