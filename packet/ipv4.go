package packet

import (
	"errors"

	"github.com/google/gopacket/layers"
)

// SetIPv4Details sets IPv4 details on inbound packets
func (p *Packet) SetIPv4Details(ipv4 *layers.IPv4) {
	p.ipv4 = ipv4
}

// IPv4Details returns the IPv4 details set on the packet
// or an error if no details are set
func (p *Packet) IPv4Details() (*layers.IPv4, error) {
	if p.ipv4 == nil {
		return nil, errors.New("no ipv4 details for packet")
	}
	return p.ipv4, nil
}
