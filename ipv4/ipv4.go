package ipv4

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"net"

	"github.com/adrianosela/rdtp"
	"github.com/adrianosela/rdtp/packet"
	"github.com/pkg/errors"
)

// IPv4 represents the underlying IPv4 network
// and functions to interact with the network
// interface
type IPv4 struct {
	sck int
}

// NewIPv4 returns a new ipv4 network interface
func NewIPv4() (*IPv4, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, rdtp.IPProtoRDTP)
	if err != nil {
		return nil, errors.Wrap(err, "could not get raw network socket")
	}
	return &IPv4{
		sck: fd,
	}, nil
}

// Send sends a packet to the destination IP address
func (ip *IPv4) Send(dstIP string, pck *packet.Packet) error {
	remoteAddr, err := parseAddr(dstIP)
	if err != nil {
		return errors.Wrap(err, "could not parse ip address")
	}
	if err := syscall.Sendto(ip.sck, pck.Serialize(), 0, remoteAddr); err != nil {
		return errors.Wrap(err, "could not send data to network socket")
	}
	return nil
}

func (ip *IPv4) PacketListener(next func(*packet.Packet) error) error {
	// readable file for socket's file descriptor
	f := os.NewFile(uintptr(ip.sck), fmt.Sprintf("fd %d", ip.sck))

	fmt.Println("listening network interfaces: /* TODO */")

	for {
		buf := make([]byte, 65535) // maximum IP packet

		ipDatagramSize, err := f.Read(buf)
		if err != nil {
			log.Println(errors.Wrap(err, "could not read data from network socket"))
			continue
		}

		rawIP := []byte(buf)[:ipDatagramSize]
		ihl := 4 * (rawIP[0] & byte(15))
		rawRDTP := rawIP[ihl:]

		rdtpPacket, err := packet.Deserialize(rawRDTP)
		if err != nil {
			log.Println(errors.Wrap(err, "could not deserialize rdtp packet"))
			continue
		}

		if err = next(rdtpPacket); err != nil {
			log.Println(errors.Wrap(err, "could not forward received rdtp packet"))
			continue
		}
	}

}

func parseAddr(ip string) (*syscall.SockaddrInet4, error) {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return nil, errors.New("invalid IPv4 address")
	}
	return &syscall.SockaddrInet4{
		Addr: [4]byte{parsed[0], parsed[1], parsed[2], parsed[3]},
	}, nil
}
