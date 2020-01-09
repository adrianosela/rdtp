package rdtp

// Packet is an RDTP packet
type Packet struct {
	SrcPort  uint8
	DstPort  uint8
	Length   uint8
	Checksum uint8
}
