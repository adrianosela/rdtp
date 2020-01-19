package rdtp

var (
	// FlagOpen is used to indicate that a host wants to start a
	// new connection with another. Is analogous to TCP's SYN flag.
	FlagOpen = uint16(65535) // 1000_0000_0000_0000
	// FlagClose is used to indicate that a host wants to terminate
	// communication with another. Is analogous to TCP's FIN flag.
	FlagClose = uint16(32767) // 0100_0000_0000_0000

//  Flag
)
