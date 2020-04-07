package packet

import (
	"errors"
	"net"
)

// SetDestinationIPv4 sets the destination IPv4 on the packet
func (p *Packet) SetDestinationIPv4(ip net.IP) {
	p.dstIP = ip
}

// SetSourceIPv4 sets the source IPv4 on the packet
func (p *Packet) SetSourceIPv4(ip net.IP) {
	p.srcIP = ip
}

// GetDestinationIPv4 returns the destination IPv4
// set on the packet or an error if none is set
func (p *Packet) GetDestinationIPv4() (net.IP, error) {
	if p.dstIP == nil {
		return nil, errors.New("no destination IPv4 address set on packet")
	}
	return p.dstIP, nil
}

// GetSourceIPv4 returns the source IPv4
// set on the packet or an error if none is set
func (p *Packet) GetSourceIPv4() (net.IP, error) {
	if p.srcIP == nil {
		return nil, errors.New("no source IPv4 address set on packet")
	}
	return p.srcIP, nil
}
