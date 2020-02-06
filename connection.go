package rdtp

import (
	"net"
)

// Conn is an RDTP connection
type Conn struct {
	lPort uint16 // local port
	rPort uint16 // remote port

	// TODO
}

// Dial establishes an RDTP connection with a remote host
func Dial(ip *net.IP, p uint16) (*Conn, error) {
	// TODO

	// local port has to be requested from
	// the rdtp controller running on this host

	return &Conn{rPort: p}, nil
}
