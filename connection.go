package rdtp

import (
	"net"

	"github.com/pkg/errors"
)

// Conn is a logical communication channel
// between the local and remote hosts.
// Implements the net.Conn interface
// https://golang.org/pkg/net/#Conn
type Conn struct {
	rdtp net.Conn
}

// Dial returns a connection to a remote address
// where the remote address has a format:
// ${host}:${port}
func Dial(address string) (*Conn, error) {
	c, err := net.Dial("unix", DefaultRDTPServiceAddr)
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to rdtp service")
	}
	// FIXME, makes this json
	if _, err := c.Write([]byte(address)); err != nil {
		return nil, errors.Wrap(err, "could not send address to rdtp service")
	}
	return &Conn{rdtp: c}, nil
}

// Read reads data from the connection.
// Read can be made to time out and return an Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetReadDeadline.
func (c *Conn) Read(b []byte) (n int, err error) {
	return c.rdtp.Read(b)
}

// Write writes data to the connection.
// Write can be made to time out and return an Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetWriteDeadline.
func (c *Conn) Write(b []byte) (n int, err error) {
	return c.rdtp.Write(b)
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (c *Conn) Close() error {
	return c.Close()
}
