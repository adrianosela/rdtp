package rdtp

import (
	"net"
	"time"

	"github.com/pkg/errors"
)

// Conn implements *part* of the net.Conn interface
// https://golang.org/pkg/net/#Conn
type Conn struct {
	client net.Conn
	laddr  net.Addr
	raddr  net.Addr
}

// Dial establishes an RDTP connection with a remote IP host
func Dial(address string) (*Conn, error) {
	c, err := net.Dial("unix", DefaultRDTPServiceAddr)
	if err != nil {
		return nil, errors.Wrap(err, "could not acquire RDTP service connection")
	}
	// laddr = network IP (i.e. 192.168.1/24) (will be changed by NAT)
	// raddr = destination IP given
	return &Conn{
		client: c,
	}, nil
}

// Close closes an RDTP connection
func (c *Conn) Close() error {
	return c.client.Close()
}

// Read reads RDTP data onto the given buffer
func (c *Conn) Read(buf []byte) (int, error) {
	return c.client.Read(buf)
}

// Write writes data to an RDTP connection
func (c *Conn) Write(data []byte) (int, error) {
	return c.client.Write(data)
}

// LocalAddr returns the local network address.
func (c *Conn) LocalAddr() net.Addr {
	return c.laddr
}

// RemoteAddr returns the remote network address.
func (c *Conn) RemoteAddr() net.Addr {
	return c.raddr
}
