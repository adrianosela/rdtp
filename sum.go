package rdtp

func computeChecksum(p *Packet) uint16 {
	csum, length := uint16(0), len(p.Payload)-1
	for i := 0; i < length; i += 2 {
		csum += uint16(p.Payload[i]) << 8
		csum += uint16(p.Payload[i+1])
	}
	if len(p.Payload)%2 == 1 {
		csum += uint16(p.Payload[length]) << 8
	}
	return ^uint16(csum)
}
