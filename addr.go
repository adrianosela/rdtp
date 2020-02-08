package rdtp

import (
	"fmt"
)

// Addr represents an RDTP address.
// Implements the net.Addr interface.
type Addr struct {
	port uint16
	ip   string
}

// Network returns the name of the network
func (a *Addr) Network() string {
	return Network
}

// String returns the string form of the address
func (a *Addr) String() string {
	return fmt.Sprintf("%s:%d", a.ip, a.port)
}
