package rdtp

import (
	"net"
	"time"

	"github.com/pkg/errors"
)

// Conn is a logical communication channel between the local and remote hosts.
// Implements the net.Conn interface (https://golang.org/pkg/net/#Conn)
type Conn struct {
	laddr *Addr
	raddr *Addr
	svc   net.Conn
}

// Dial returns a connection to a remote address
// where the remote address has a format: ${host}:${port}
func Dial(address string) (*Conn, error) {
	svc, err := net.Dial("unix", DefaultRDTPServiceAddr)
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to rdtp service")
	}

	raddr, err := fromString(address)
	if err != nil {
		return nil, errors.Wrap(err, "invalid remote rdtp address")
	}

	req, err := NewClientMessage(ClientMessageTypeDial, nil, raddr)
	if err != nil {
		return nil, errors.Wrap(err, "could not create rdtp dial request")
	}

	// first message out must be the remote rdtp address (e.g. ${host}:${port})
	if _, err := svc.Write(req); err != nil {
		return nil, errors.Wrap(err, "could not send address to rdtp service")
	}

	verifiedLocalAddr, err := waitForServiceMessageOK(svc)
	if err != nil {
		return nil, errors.Wrap(err, "could not receive OK message from service")
	}

	return &Conn{
		svc:   svc,
		laddr: verifiedLocalAddr,
		raddr: raddr,
	}, nil
}

// Read reads data from the connection.
func (c Conn) Read(b []byte) (n int, err error) {
	return c.svc.Read(b)
}

// Write writes data to the connection.
func (c Conn) Write(b []byte) (n int, err error) {
	return c.svc.Write(b)
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (c Conn) Close() error {
	return c.svc.Close()
}

// LocalAddr returns the local address for this conn
func (c Conn) LocalAddr() net.Addr {
	return c.laddr
}

// RemoteAddr returns the remote address for this conn
func (c Conn) RemoteAddr() net.Addr {
	return c.raddr
}

// SetDeadline sets the deadline on the connection
func (c Conn) SetDeadline(t time.Time) error {
	return c.svc.SetDeadline(t)
}

// SetReadDeadline sets the read deadline on the connection
func (c Conn) SetReadDeadline(t time.Time) error {
	return c.svc.SetReadDeadline(t)

}

// SetWriteDeadline sets the write deadline on the connection
func (c Conn) SetWriteDeadline(t time.Time) error {
	return c.svc.SetWriteDeadline(t)
}
