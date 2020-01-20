package controller

import (
	"net"

	"github.com/pkg/errors"
)

// Worker handles the local transport layer processing
// for a single process-process communication
type Worker struct {
	txConn net.Conn
	txPort uint16
}

// NewWorker returns an RDTP Worker struct
func NewWorker(dstIP string, dstPort uint16) (*Worker, error) {
	// resolve destination address
	dst, err := net.ResolveIPAddr("ip", dstIP)
	if err != nil {
		return nil, errors.Wrap(err, "could not resolve destination IP address")
	}
	txIPConn, err := net.DialIP("ip:ip", nil, dst)
	if err != nil {
		return nil, errors.Wrap(err, "could not dial IP")
	}
	// build connection object
	return &Worker{
		txConn: txIPConn,
		txPort: dstPort,
	}, nil
}

// Close gracefully shuts down a worker
func (c *Worker) Close() error {
	// TODO
	return nil
}
