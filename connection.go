package rdtp

import (
	"github.com/adrianosela/rdtp/worker"
	"github.com/pkg/errors"
)

// Conn is an RDTP connection
type Conn struct {
	worker *worker.Worker
}

// Dial establishes an RDTP connection with a remote IP host
func Dial(ip string) (*Conn, error) {
	// define destination address (on loopback for now)
	w, err := worker.NewWorker(ip)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start rdtp worker")
	}
	return &Conn{worker: w}, nil
}

// Close closes an RDTP connection
func (c *Conn) Close() error {
	if err := c.worker.Kill(); err != nil {
		return errors.Wrap(err, "failed to terminate rdtp worker gracefully")
	}
	return nil
}
