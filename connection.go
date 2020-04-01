package rdtp

import (
	"net"
	"time"

	"github.com/pkg/errors"
)

// Conn implements the net.Conn interface
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

// SetDeadline sets the read and write deadlines associated
// with the connection.
func (c *Conn) SetDeadline(t time.Time) error {
	return c.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.SetWriteDeadline(t)
}
