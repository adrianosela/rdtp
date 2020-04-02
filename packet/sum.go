package packet

// SetSum sets the checksum on an rdtp packet
func (p *Packet) SetSum() {
	p.Checksum = p.sum()
}

// CheckSum verifies the checksum on an rdtp packet
func (p *Packet) CheckSum() bool {
	return p.Checksum == p.sum()
}

func (p *Packet) sum() uint16 {
	var csum uint16

	csum += p.SrcPort
	csum += p.DstPort
	csum += p.Length

	csum += uint16(p.SeqNo >> 8)
	csum += uint16(p.SeqNo)
	csum += uint16(p.AckNo >> 8)
	csum += uint16(p.AckNo)

	csum += uint16(p.Flags)

	for i := 0; i < len(p.Payload); i++ {
		csum += uint16(p.Payload[i])
	}

	return ^csum
}
