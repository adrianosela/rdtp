package rdtp

import "errors"

// Listener is an RDTP listner
type Listener struct {
	addr *Addr
	// TODO
}

// Listen listens for new RDTP connections on a given port
func Listen(p uint16) (*Listener, error) {
	return &Listener{
		addr: &Addr{port: p},
		// TODO
	}, nil
}

// Accept accepts a new connection on a given listener
func (l *Listener) Accept() (*Conn, error) {
	return &Conn{
		// TODO
	}, nil
}

// Close closes all active connections for a given listener
func (l *Listener) Close() error {
	// TODO
	return nil
}

// Addr returns the network address of a listener
func (l *Listener) Addr() (*Addr, error) {
	if l.addr == nil {
		return nil, errors.New("no address associated with listener")
	}
	return l.addr, nil
}
