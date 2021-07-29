package rdtp

import (
	"io"
	"net"

	"github.com/pkg/errors"
)

// Listener listens for new rdtp
// connections on the interface address.
// Implements the net.Listener interface
// https://golang.org/pkg/net/#Listener
type Listener struct {
	laddr *Addr
	svc   net.Conn
}

// Listen announces on the local network address
func Listen(address string) (net.Listener, error) {
	svc, err := net.Dial("unix", DefaultRDTPServiceAddr)
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to rdtp service")
	}

	laddr, err := fromString(address)
	if err != nil {
		return nil, errors.Wrap(err, "address is not a valid rdtp address")
	}

	req, err := NewRequest(RequestTypeListen, laddr, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create listen request for rdtp service")
	}

	if _, err = svc.Write(req); err != nil {
		return nil, errors.Wrap(err, "could not send listen request to rdtp service")
	}

	l := &Listener{
		laddr: laddr,
		svc:   svc,
	}

	return l, nil
}

// Accept waits for and returns the next connection to the listener.
func (l *Listener) Accept() (net.Conn, error) {
	buf := make([]byte, 1024)
	n, err := l.svc.Read(buf)
	if err != nil {
		if err == io.EOF {
			// TODO: handle conn closed by rdtp service
		}
		return nil, errors.Wrap(err, "could not read remote address notification")
	}

	raddr, err := fromString(string(buf[:n]))
	if err != nil {
		return nil, errors.Wrap(err, "remote address is not valid")
	}

	svc, err := net.Dial("unix", DefaultRDTPServiceAddr)
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to rdtp service")
	}

	req, err := NewRequest(RequestTypeAccept, l.laddr, raddr)
	if err != nil {
		return nil, errors.Wrap(err, "could not create accept request for rdtp service")
	}

	if _, err = svc.Write(req); err != nil {
		return nil, errors.Wrap(err, "could not send accept request to rdtp service")
	}

	return &Conn{
		laddr: l.laddr,
		raddr: raddr,
		svc:   svc,
	}, nil
}

// Close closes the listener.
func (l *Listener) Close() error {
	return l.svc.Close()
}

// Addr returns the listener's network address.
func (l *Listener) Addr() net.Addr {
	return l.laddr
}
