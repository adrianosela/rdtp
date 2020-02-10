package rdtp

import (
	"net"

	"github.com/pkg/errors"
)

// Conn is an RDTP connection
type Conn struct {
	client net.Conn
}

// Dial establishes an RDTP connection with a remote IP host
func Dial(ip []byte) (*Conn, error) {
	c, err := net.Dial("unix", DefaultRDTPServiceAddr)
	if err != nil {
		return nil, errors.Wrap(err, "could not acquire RDTP service connection")
	}
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
