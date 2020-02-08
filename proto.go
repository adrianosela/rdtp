package rdtp

const (
	// IPProtoRDTP is the protocol number for RDTP packets over IP
	// The value 157 (0x9D) is unassigned as per:
	// https://en.wikipedia.org/wiki/List_of_IP_protocol_numbers
	IPProtoRDTP = 0x9D

	// Network is the name of the RDTP network
	Network = "rdtp"

	// DiscoveryPort is the port that receives SYN packets
	DiscoveryPort = uint16(0)
)
