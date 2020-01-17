package rdtp

func computeChecksum(p *Packet) uint16 {
	csum := uint16(0)
	csum += p.SrcPort
	csum += p.DstPort
	csum += p.Length
	for i := 0; i < len(p.Payload); i++ {
		csum += uint16(p.Payload[i])
	}
	return ^csum
}

func verifyChecksum(p *Packet) bool {
	return p.Checksum == computeChecksum(p)
}
