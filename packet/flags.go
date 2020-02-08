package packet

const (
	synMask = 0x80
	ackMask = 0x40
	finMask = 0x20
	errMask = 0x10
)

// SetSYN sets the SYN flag on a packet
func (p *Packet) SetSYN() {
	p.Flags = p.Flags | synMask
}

// SetACK sets the ACK flag on a packet
func (p *Packet) SetACK() {
	p.Flags = p.Flags | ackMask
}

// SetFIN sets the FIN flag on a packet
func (p *Packet) SetFIN() {
	p.Flags = p.Flags | finMask
}

// SetERR sets the ERR flag on a packet
func (p *Packet) SetERR() {
	p.Flags = p.Flags | errMask
}

// IsSYN returns true if the SYN flag is set
func (p *Packet) IsSYN() bool {
	return p.Flags&synMask != 0
}

// IsACK returns true if the ACK flag is set
func (p *Packet) IsACK() bool {
	return p.Flags&ackMask != 0
}

// IsFIN returns true if the FIN flag is set
func (p *Packet) IsFIN() bool {
	return p.Flags&finMask != 0
}

// IsERR returns true if the ERR flag is set
func (p *Packet) IsERR() bool {
	return p.Flags&errMask != 0
}
