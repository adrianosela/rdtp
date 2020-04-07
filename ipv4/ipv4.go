package ipv4

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/adrianosela/rdtp"
	"github.com/adrianosela/rdtp/packet"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/pkg/errors"
)

// IPv4 represents the underlying IPv4 network
// and functions to interact with the network
// interface
type IPv4 struct {
	sckfd int
}

// NewIPv4 returns a new ipv4 network interface
func NewIPv4() (*IPv4, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, rdtp.IPProtoRDTP)
	if err != nil {
		return nil, errors.Wrap(err, "could not get raw network socket")
	}
	return &IPv4{
		sckfd: fd,
	}, nil
}

// Send sends a packet to the destination IP address
func (ip *IPv4) Send(pck *packet.Packet) error {
	dstIP, err := pck.GetDestinationIPv4()
	if err != nil {
		return errors.Wrap(err, "could not determine destination IP addresss")
	}

	remote := &syscall.SockaddrInet4{
		Addr: [4]byte{dstIP[0], dstIP[1], dstIP[2], dstIP[3]},
	}

	if err := syscall.Sendto(ip.sckfd, pck.Serialize(), 0, remote); err != nil {
		return errors.Wrap(err, "could not send data to network socket")
	}
	return nil
}

// Receive forwards all ipv4 packets received which carry rdtp
// These ipv4 packets are processed until an rdtp packet.Packet
// is extracted, then the next() function is called with the packet
func (ip *IPv4) Receive(next func(*packet.Packet) error) error {
	rdtpFile := os.NewFile(uintptr(ip.sckfd), fmt.Sprintf("fd %d", ip.sckfd))
	buf := make([]byte, 65535) // maximum IP packet
	for {
		ipDatagramSize, err := rdtpFile.Read(buf)
		if err != nil {
			log.Println(errors.Wrap(err, "could not read data from network socket"))
			continue
		}

		// decode ipv4
		networkPck := gopacket.NewPacket(
			buf[:ipDatagramSize],
			layers.LayerTypeIPv4,
			gopacket.Default)
		ipv4NetworkData := networkPck.Layer(layers.LayerTypeIPv4)

		if ipv4NetworkData == nil {
			log.Println("not an ipv4 packet")
			continue
		}

		ipv4 := ipv4NetworkData.(*layers.IPv4)

		rdtpPacket, err := packet.Deserialize(ipv4.Payload)
		if err != nil {
			log.Println(errors.Wrap(err, "could not deserialize rdtp packet"))
			continue
		}

		rdtpPacket.SetDestinationIPv4(ipv4.DstIP)
		rdtpPacket.SetSourceIPv4(ipv4.SrcIP)

		if err = next(rdtpPacket); err != nil {
			log.Println(errors.Wrap(err, "could not process received rdtp packet"))
			continue
		}
	}
}
