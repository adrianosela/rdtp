package rdtp

import (
	"net"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Listener implements the net.Listener interface
// https://golang.org/pkg/net/#Listener
type Listener struct {
	addr *Addr
}

// Listen returns an rdtp listener
func Listen(address string) (*Listener, error) {
	addr := strings.Split(address, ":")

	host := addr[0]
	port, err := strconv.ParseUint(addr[1], 10, 16)
	if err != nil {
		return nil, errors.Wrap(err, "invalid port number")
	}

	// TODO: talk to the rdtp service to reserve port
	//

	return &Listener{
		addr: &Addr{
			ip:   host,
			port: uint16(port),
		},
	}, nil
}

// Accept waits for and returns the next connection to the listener.
func (l *Listener) Accept() (net.Conn, error) {
	// TODO:
	return nil, nil
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l *Listener) Close() error {
	// TODO:
	return nil
}

// Addr returns the listener's network address.
func (l *Listener) Addr() net.Addr {
	return l.addr
}
