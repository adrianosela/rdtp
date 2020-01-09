package rdtp

const (
	// MaxPacketBytes is the maximum size of an RDTP packet incl. header
	MaxPacketBytes = 65516 // 65536 - IPv4 Header (20 bytes)
	// MinHeaderBytes is the byte size of an RDTP header
	MinHeaderBytes = 8
)
