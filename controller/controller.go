package controller

import (
	"fmt"
	"log"
	"os"
	"sync"
	"syscall"

	"github.com/adrianosela/rdtp"
	"github.com/pkg/errors"
)

const (
	// maxPortNo is the amount of ports that fit in a 16 bit domain
	// such as the 16 bit source/destination port identifiers in an RDTP packet
	maxPortNo = uint16(65535)
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
		Ports: make(map[uint16]*Worker, maxPortNo),
	}
}

// Start starts the RDTP Controller service
func (ctrl *Controller) Start() error {
	// get raw network socket (AF_INET = IPv4)
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, rdtp.IPPROTO_RDTP)
	if err != nil {
		return errors.Wrap(err, "could not get raw network socket")
	}

	// readable file for socket's file descriptor
	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))

	fmt.Println("listening on all local IPv4 network interfaces")
	for {
		buf := make([]byte, 65535) // maximum IP packet

		ipDatagramSize, err := f.Read(buf)
		if err != nil {
			log.Println(errors.Wrap(err, "could not read data from network socket"))
			continue
		}

		rawIP := []byte(buf)[:ipDatagramSize]
		ihl := 4 * (rawIP[0] & byte(15))
		rawRDTP := rawIP[ihl:]

		rdtpPacket, err := rdtp.Deserialize(rawRDTP)
		if err != nil {
			log.Println(errors.Wrap(err, "could not deserialize rdtp packet"))
			continue
		}

		if !rdtpPacket.Check() {
			log.Println("failed checksum, packet dropped")
			continue
		}

		if err = ctrl.MultiplexPacket(rdtpPacket); err != nil {
			log.Println(errors.Wrap(err, "could not multiplex rdtp packet"))
			continue
		}
	}
}

// Shutdown force-closes all existing connections for a controller
func (ctrl *Controller) Shutdown() {
	ctrl.Lock()
	defer ctrl.Unlock()

	for p, c := range ctrl.Ports {
		if err := c.Kill(); err != nil {
			log.Println(errors.Wrapf(err, "error closing rdtp conn on port %d", p))
		}
		delete(ctrl.Ports, p)
	}
}
