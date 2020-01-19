package rdtp

import (
	"fmt"
	"sync"
)

var ctrl = NewController() // global rdtp controller

const (
	// MaxPortNo is the amount of ports that fit in a 16 bit domain
	// such as the 16 bit source/destination port identifiers in an RDTP packet
	MaxPortNo = uint16(65535)
)

// Controller manages port numbers for individual RDTP connections
type Controller struct {
	sync.Mutex // inherit mutex lock behavior
	Ports      map[uint16]*Conn
}

// NewController is the constructor for a new RDTP communication controller.
// This kind of thing typically runs in the Kernel to manage ports for TCP/UDP
func NewController() *Controller {
	return &Controller{
		Ports: make(map[uint16]*Conn, MaxPortNo),
	}
}

// Allocate associates a connection with an RDTP port
func (ctrl *Controller) Allocate(c *Conn) error {
	ctrl.Lock()
	defer ctrl.Unlock()

	for port := uint16(0); port < MaxPortNo; port++ {
		if _, ok := ctrl.Ports[port]; !ok {
			ctrl.Ports[port] = c // reserve port for conn
			c.rxPort = port      // set port on conn
			return nil
		}
	}
	return fmt.Errorf("all ports in use")
}

// Deallocate frees up an RDTP port
func (ctrl *Controller) Deallocate(port uint16) {
	ctrl.Lock()
	defer ctrl.Unlock()

	delete(ctrl.Ports, port)
}

// Shutdown force-closes all existing connections for a controller
func (ctrl *Controller) Shutdown() {
	ctrl.Lock()
	defer ctrl.Unlock()

	for p, c := range ctrl.Ports {
		c.Close()
		delete(ctrl.Ports, p)
	}
}
