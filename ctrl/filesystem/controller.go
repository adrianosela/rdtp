package filesystem

import (
	"fmt"
	"time"

	"github.com/adrianosela/rdtp"
	"github.com/gofrs/flock"
	"github.com/pkg/errors"
)

// FSController is an RDTP communication controller implemented
// with the host's file system
type FSController struct {
	lock *flock.Flock
	path string
}

// NewFSController is the constructor for a new RDTP file system controller
func NewFSController(path string) (*FSController, error) {
	if path == "" {
		path = rdtpFilePath
	}

	// if no statefile at given path, create it and initialize
	// it by allocating port 0 to this controller
	if !fileExists(path) {
		s := &Statefile{
			Ports: map[uint16]int64{uint16(0): time.Now().UnixNano()},
		}
		if err := s.commit(path); err != nil {
			return nil, errors.Wrap(err, "could not initialize statefile in filesystem")
		}
	}

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

	state, err := getState(c.path)
	if err != nil {
		return 0, errors.Wrap(err, "could not get rdtp state")
	}

	for port := uint16(1); port < rdtp.MaxPortNo; port++ {
		// give out first unused port
		if _, ok := state.Ports[port]; !ok {
			state.Ports[port] = time.Now().UnixNano()
			if err := state.commit(c.path); err != nil {
				return 0, errors.Wrap(err, "could not commit statefile to filesystem")
			}
			return port, nil
		}
	}

	return 0, fmt.Errorf("all ports in use")
}

// Allocate allocates a specific RDTP port
func (c *FSController) Allocate(p uint16) error {
	if err := c.lock.Lock(); err != nil {
		return errors.Wrap(err, "could not acquire lock for rdtp file")
	}
	defer c.lock.Unlock()

	if p == 0 {
		return errors.New("port 0 cannot be used (reserved)")
	}

	state, err := getState(c.path)
	if err != nil {
		return errors.Wrap(err, "could not get rdtp state")
	}

	if _, ok := state.Ports[p]; !ok {
		state.Ports[p] = time.Now().UnixNano()
		if err := state.commit(c.path); err != nil {
			return errors.Wrap(err, "could not commit statefile to filesystem")
		}
		return nil
	}

	return errors.New("port is in use")
}

// Deallocate frees up a given RDTP port
func (c *FSController) Deallocate(p uint16) error {
	if err := c.lock.Lock(); err != nil {
		return errors.Wrap(err, "could not acquire lock for rdtp file")
	}
	defer c.lock.Unlock()

	if p == 0 {
		return errors.New("port 0 cannot be used (reserved)")
	}

	state, err := getState(c.path)
	if err != nil {
		return errors.Wrap(err, "could not get rdtp state")
	}
	delete(state.Ports, p)
	if err := state.commit(c.path); err != nil {
		return errors.Wrap(err, "could not commit statefile to filesystem")
	}

	return nil
}
