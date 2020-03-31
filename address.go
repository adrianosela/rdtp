package rdtp

import "fmt"

const (
	// Network is the name of the RDTP network
	Network = "rdtp"
)

// Addr implements the net.Addr interface
type Addr struct {
	ip   string
	port uint16
}

// Network returns the name of the network
func (a *Addr) Network() string {
	return Network
}

// String returns the string form of the address
func (a *Addr) String() string {
	return fmt.Sprintf("%s:%d", a.ip, a.port)
}
