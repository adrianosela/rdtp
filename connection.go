package rdtp

import (
	"net"

	"github.com/pkg/errors"
)

// Conn handles the local transport layer processing
// for a single process-process communication
type Conn struct {
	txConn net.Conn
	txPort uint16
}

// NewConn returns an RDTP Connection struct
func NewConn(dstIP string, dstPort uint16) (*Conn, error) {
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
	return &Conn{
		txConn: txIPConn,
		txPort: dstPort,
	}, nil
}

func (c *Conn) Close() error {
	// TODO
	return nil
}
