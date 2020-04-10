package rdtp

import (
	"net"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Listener listens for new rdtp
// connections on the interface address.
// Implements the net.Listener interface
// https://golang.org/pkg/net/#Listener
type Listener struct {
	addr *Addr
}

// Listen announces on the local network address
func Listen(address string) (net.Listener, error) {
	splt := strings.Split(address, ":")

	if len(splt) > 2 || len(splt) == 0 {
		return nil, errors.New("invalid ipv4 address")
	}

	var a *Addr

	// if no port given
	if len(splt) == 1 {
		a.Host = splt[0]
		a.Port = DiscoveryPort
	}

	// if port is given (and host may or may not be)
	if len(splt) <= 2 {
		a.Host = splt[0] // host might be empty, which is okay
		if splt[1] == "" {
			a.Port = uint16(0)
		} else {
			port, err := strconv.ParseUint(splt[1], 10, 16)
			if err != nil {
				return nil, errors.Wrap(err, "invalid port number")
			}
			a.Port = uint16(port)
		}
	}

	// TODO
	return &Listener{addr: a}, nil
}

// Accept waits for and returns the next connection to the listener.
func (l *Listener) Accept() (net.Conn, error) {
	// TODO
	return nil, nil
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l *Listener) Close() error {
	// TODO
	return nil
}

// Addr returns the listener's network address.
func (l *Listener) Addr() net.Addr {
	return l.addr
}
