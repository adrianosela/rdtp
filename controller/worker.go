package controller

import (
	"net"

	"github.com/pkg/errors"
)

// Worker handles the local transport layer processing
// for a single process-process communication
type Worker struct {
	// TODO
}

// NewWorker returns an RDTP Worker struct
func NewWorker() (*Worker, error) {
	return &Worker{
		// TODO
	}, nil
}

// Kill shuts down a worker
func (c *Worker) Kill() error {
	// TODO
	return nil
}
