package rdtp

// Conn uniquely identifies a logical communication
// channel between the local and remote hosts
type Conn struct {
	LAddr *Addr
	Raddr *Addr
}
