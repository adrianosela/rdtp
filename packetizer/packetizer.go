package packetizer

// Packetizer chunks data into MSS bytes and
// wraps them into packets
type Packetizer interface {
	Wrap([]byte, *rdtp.Conn, func(*packet.Packet) error)
}
