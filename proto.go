package rdtp

const (
	// IPProtoRDTP is the protocol number for RDTP packets over IP
	// The value 157 (0x9D) is unassigned as per:
	// https://en.wikipedia.org/wiki/List_of_IP_protocol_numbers
	IPProtoRDTP = 0x9D

	// DiscoveryPort is the port that receives new RDTP connections
	DiscoveryPort = uint16(0)

	// MaxPort is the highest possible RDTP port number
	MaxPort = uint16(65535)

	// DefaultRDTPServiceAddr is the default rdtp service socket
	DefaultRDTPServiceAddr = "/var/run/rdtp.sock"
)
