package controller

import (
	"fmt"
)

// Allocate associates a connection with an RDTP port
func (ctrl *Controller) Allocate(c *Worker) (uint16, error) {
	ctrl.Lock()
	defer ctrl.Unlock()

	for port := uint16(0); port < MaxPortNo; port++ {
		if _, ok := ctrl.Ports[port]; !ok {
			ctrl.Ports[port] = c // reserve port for conn
			return port, nil
		}
	}
	return 0, fmt.Errorf("all ports in use")
}

// Deallocate frees up an RDTP port
func (ctrl *Controller) Deallocate(port uint16) {
	ctrl.Lock()
	defer ctrl.Unlock()

	delete(ctrl.Ports, port)
}
