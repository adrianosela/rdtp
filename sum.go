package rdtp

// this checksum can be improved
func computeChecksum(p *Packet) uint16 {
	csum := uint16(0)
	for i := 0; i < len(p.Payload); i++ {
		csum += uint16(p.Payload[i])
	}
	return uint16(csum)
}
