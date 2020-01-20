package controller

import (
	"log"
	"sync"

	"github.com/pkg/errors"
)

const (
	// MaxPortNo is the amount of ports that fit in a 16 bit domain
	// such as the 16 bit source/destination port identifiers in an RDTP packet
	MaxPortNo = uint16(65535)

	// the ip header is 20 bytes
	ipHeaderLen = 20
)

// Controller is the RDTP communication controller.
// This kind of thing typically runs in the Kernel to manage ports for TCP/UDP
type Controller struct {
	sync.Mutex // inherit mutex lock behavior
	Ports      map[uint16]*Worker
}

// NewController is the constructor for a new RDTP communication controller.
// This kind of thing typically runs in the Kernel to manage ports for TCP/UDP
func NewController() *Controller {
	return &Controller{
		Ports: make(map[uint16]*Worker, MaxPortNo),
	}
}

// Start starts the RDTP Controller service
func (ctrl *Controller) Start(network string) error {
	switch network {
	case "ip4":
		return ctrl.listenRDTPOverIPv4()
	}
	return errors.New("network not supported. supported networks: \"ip4\"")
}

// Shutdown force-closes all existing connections for a controller
func (ctrl *Controller) Shutdown() {
	ctrl.Lock()
	defer ctrl.Unlock()

	for p, c := range ctrl.Ports {
		if err := c.Close(); err != nil {
			log.Println(errors.Wrapf(err, "error closing rdtp conn on port %d", p))
		}
		delete(ctrl.Ports, p)
	}
}
