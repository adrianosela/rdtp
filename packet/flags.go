package packet

const (
	synMask = 0x80
	ackMask = 0x40
	finMask = 0x20
	errMask = 0x10
)

// SetFlagSYN sets the SYN flag on a packet
func (p *Packet) SetFlagSYN() {
	p.Flags = p.Flags | synMask
}

// SetFlagACK sets the ACK flag on a packet
func (p *Packet) SetFlagACK() {
	p.Flags = p.Flags | ackMask
}

// SetFlagFIN sets the FIN flag on a packet
func (p *Packet) SetFlagFIN() {
	p.Flags = p.Flags | finMask
}

// SetFlagERR sets the ERR flag on a packet
func (p *Packet) SetFlagERR() {
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
