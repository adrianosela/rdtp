package filesystem

import (
	"fmt"

	"github.com/gofrs/flock"
	"github.com/pkg/errors"
)

const file = "~/.rdtp"

// FSController is an RDTP communication controller implemented
// with the host's file system
type FSController struct {
	lock *flock.Flock
	path string
}

// NewFSController is the constructor for a new RDTP file system controller
func NewFSController(path string) (*FSController, error) {
	return &FSController{
		lock: flock.New(path),
		path: path,
	}, nil
}

// AllocateAny allocates an available RDTP port
func (c *FSController) AllocateAny() (uint16, error) {
	if err := c.lock.Lock(); err != nil {
		return 0, errors.Wrap(err, "could not acquire lock for rdtp file")
	}
	defer c.lock.Unlock()

	// TODO

	return 0, fmt.Errorf("all ports in use")
}

// Allocate allocates a specific RDTP port
func (c *FSController) Allocate(p uint16) error {
	if err := c.lock.Lock(); err != nil {
		return errors.Wrap(err, "could not acquire lock for rdtp file")
	}
	defer c.lock.Unlock()

	// TODO

	return nil
}

// Deallocate frees up a given RDTP port
func (c *FSController) Deallocate(port uint16) error {
	if err := c.lock.Lock(); err != nil {
		return errors.Wrap(err, "could not acquire lock for rdtp file")
	}
	defer c.lock.Unlock()

	// TODO

	return nil
}
