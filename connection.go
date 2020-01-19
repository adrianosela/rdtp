package rdtp

import (
	"net"

	"github.com/pkg/errors"
)

// Conn handles the local transport layer processing
// for a single process-process communication
type Conn struct {
	// lower layer communication
	// channel pointers
	rxConn net.Conn
	txConn net.Conn
	// process identifyers
	// locally and remotely
	rxPort uint16
	txPort uint16
}

// DialRDTP returns an RDTP Connection struct
func DialRDTP(dstIP string, dstPort uint16) (*Conn, error) {
	// resolve listener address
	src, err := net.ResolveIPAddr("ip", "127.0.0.1")
	if err != nil {
		return nil, errors.Wrap(err, "could not resolve source IP address")
	}
	rxIPConn, err := net.ListenIP("ip:ip", src)
	if err != nil {
		return nil, errors.Wrap(err, "could not listen for IP")
	}
	// resolve destination address
	dst, err := net.ResolveIPAddr("ip", dstIP)
	if err != nil {
		return nil, errors.Wrap(err, "could not resolve destination IP address")
	}
	txIPConn, err := net.DialIP("ip:ip", nil, dst)
	if err != nil {
		return nil, errors.Wrap(err, "could not dial IP")
	}

	c := &Conn{
		rxConn: rxIPConn,
		txConn: txIPConn,
		txPort: dstPort,
	}

	if err := ctrl.Allocate(c); err != nil {
		return nil, errors.Wrap(err, "could not allocate RDTP port for connection")
	}

	return c, nil
}

// func (c *Conn) Read(b []byte) (n int, err error) {
// 	// allocate buffer to read as large a packet as possible
// 	buf := make([]byte, MaxPacketBytes)
//
// 	ipPayloadSize, err := c.channel.Read(buf)
// 	if err != nil {
// 		return 0, errors.Wrap(err, "could not read rdtp packet from underlying datagram")
// 	}
//
// 	p, err := Deserialize([]byte(buf)[c.llHeaderSize:ipPayloadSize])
// 	if err != nil {
// 		return 0, errors.Wrap(err, "could not build received rdtp packet")
// 	}
//
// 	copy(b, p.Payload)
//
// 	return len(p.Payload), nil
// }
//
// // Write implements the net.Conn Write method, yet it has some limitations.
// // - it will error out if the data length is larger than the max payload size
// //   i.e. chunking data into packets is not supported at the moment.
// func (c *Conn) Write(b []byte) (n int, err error) {
// 	p, err := NewPacket(c.srcPort, c.dstPort, b)
// 	if err != nil {
// 		return 0, errors.Wrap(err, "could not build rdtp packet for sending")
// 	}
// 	sent, err := c.channel.Write(p.Serialize())
// 	if err != nil {
// 		return sent, errors.Wrap(err, "could not write rdtp packet to underlying datagram")
// 	}
// 	if sent < len(b) {
// 		return sent, errors.Wrap(err, fmt.Sprintf("only %d out of %d bytes sent", sent, len(b)))
// 	}
// 	return sent, nil
// }

// Close terminates a connection and frees up all associated resources
func (c *Conn) Close() {
	// TODO
}
