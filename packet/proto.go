package packet

const (
	// MaxPacketBytes is the maximum size of an RDTP packet incl. header
	MaxPacketBytes = 1500 // will chunk otherwise
	// HeaderByteSize is the byte size of an RDTP header
	HeaderByteSize = 9
	// MaxPayloadByteSize is the maximum size of a payload that a single RDTP
	// packet can carry
	MaxPayloadByteSize = MaxPacketBytes - HeaderByteSize
)
