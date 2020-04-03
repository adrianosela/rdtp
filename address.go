package rdtp

import "fmt"

// Network is the name of the RDTP network
const Network = "rdtp"

// Port is an abstraction for rdtp port number
type Port uint16

// Addr implements the net.Addr interface
// https://golang.org/pkg/net/#Addr
type Addr struct {
	Host string `json:"host"`
	Port Port   `json:"port"`
}

// Network returns the name of the network
func (a *Addr) Network() string {
	return Network
}

// String returns the string form of the address
func (a *Addr) String() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}
