package rdtp

// Conn is an RDTP connection
type Conn struct {
	// TODO
}

// Dial establishes an RDTP connection with a remote IP host
func Dial(ip string) (*Conn, error) {
	return &Conn{ /* TODO */ }, nil
}

// Close closes an RDTP connection
func (c *Conn) Close() error {
	// TODO
	return nil
}
