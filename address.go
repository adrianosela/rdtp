package rdtp

import "fmt"

// Network is the name of the RDTP network
const Network = "rdtp"

// Addr implements the net.Addr interface
// https://golang.org/pkg/net/#Addr
type Addr struct {
	IP   string `json:"ip"`
	Port uint16 `json:"port"`
}

// Network returns the name of the network
func (a *Addr) Network() string {
	return Network
}

// String returns the string form of the address
func (a *Addr) String() string {
	return fmt.Sprintf("%s:%d", a.IP, a.Port)
}
