package socket

import "net"

// Listener represents a connection to the application listening on a given port
type Listener struct {
	port     uint16
	notifyTo net.Conn
}

// NewListener is the Listener constructor
func NewListener(port uint16, c net.Conn) *Listener {
	return &Listener{
		port:     port,
		notifyTo: c,
	}
}
