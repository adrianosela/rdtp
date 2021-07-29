package rdtp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Network is the name of the RDTP network
const Network = "rdtp"

// Addr implements the net.Addr interface
// https://golang.org/pkg/net/#Addr
type Addr struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

// Network returns the name of the network
func (a *Addr) Network() string {
	return Network
}

// String returns the string form of the address
func (a *Addr) String() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

func fromString(address string) (*Addr, error) {
	addrParts := strings.Split(address, ":")
	if len(addrParts) > 2 || len(addrParts) == 0 {
		return nil, errors.New("invalid ipv4 address")
	}

	var a Addr

	// if no port given
	if len(addrParts) == 1 {
		a.Host = addrParts[0]
		a.Port = DiscoveryPort
	}

	// if port is given
	if len(addrParts) <= 2 {
		a.Host = addrParts[0] // host might be empty, which is okay
		if addrParts[1] == "" {
			a.Port = uint16(0)
		} else {
			port, err := strconv.ParseUint(addrParts[1], 10, 16)
			if err != nil {
				return nil, errors.Wrap(err, "invalid port number")
			}
			a.Port = uint16(port)
		}
	}

	return &a, nil
}
