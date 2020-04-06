package ipv4

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"net"

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
func (ip *IPv4) Send(dstIP string, pck *packet.Packet) error {
	remoteAddr, err := parseAddr(dstIP)
	if err != nil {
		return errors.Wrap(err, "could not parse ip address")
	}
	if err := syscall.Sendto(ip.sckfd, pck.Serialize(), 0, remoteAddr); err != nil {
		return errors.Wrap(err, "could not send data to network socket")
	}
	return nil
}

// ForwardRDTP forwards all received IP packets carrying rdtp
func (ip *IPv4) ForwardRDTP(fw func(*packet.Packet) error) error {
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

		rdtpPacket.SetIPv4(ipv4)
		if err = fw(rdtpPacket); err != nil {
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
