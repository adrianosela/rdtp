package rdtp

const (
	// MaxPacketBytes is the maximum size of an RDTP packet incl. header
	MaxPacketBytes = 65515 // 65535 - IPv4 Header (20 bytes)
	// HeaderByteSize is the byte size of an RDTP header
	HeaderByteSize = 8
	// MaxPayloadByteSize is the maximum size of a payload that a single RDTP
	// packet can carry
	MaxPayloadByteSize = MaxPacketBytes - HeaderByteSize

	// IPPROTO_RDTP is the protocol number for RDTP packets over IP
	// The value 157 (0x9D) is unassigned as per:
	// https://en.wikipedia.org/wiki/List_of_IP_protocol_numbers
	IPPROTO_RDTP = 0x9D
)
