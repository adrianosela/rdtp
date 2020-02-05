package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"encoding/binary"

	"github.com/adrianosela/rdtp"
	"github.com/pkg/errors"
)

func main() {
	// get raw network socket (AF_INET = IPv4)
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, rdtp.IPPROTO_RDTP)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not get raw network socket"))
	}
	// readable file for socket's file descriptor
	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))

	fmt.Println("listening on all local network interfaces")
	for {
		buf := make([]byte, 65535) // maximum IP packet

		ipDatagramSize, err := f.Read(buf)
		if err != nil {
			log.Println(errors.Wrap(err, "could not read data from network socket"))
			continue
		}

		rawIP := []byte(buf)[:ipDatagramSize]
		ipHeader, rawRDTP := rawIP[:20], rawIP[20:]
		rdtpHeader := rawRDTP[:rdtp.HeaderByteSize]

		fmt.Printf("IP HEADER: %v\n", ipHeader)
		printIPHeader(ipHeader)

		rdtpPacket, err := rdtp.Deserialize(rawRDTP)
		if err != nil {
			log.Println(errors.Wrap(err, "could not deserialize rdtp packet"))
			continue
		}

		fmt.Printf("RDTP HEADER: %v\n", rdtpHeader)
		printRDTPHeader(rdtpPacket)

		fmt.Printf("RDTP PAYLOAD: %s\n", string(rdtpPacket.Payload))
	}

}

func printIPHeader(header []byte) error {
	if len(header) < 20 {
		return errors.New("too small")
	}
	// 4 bit version
	version := uint8(header[0] >> 4)
	// 4 bit ihl (# of 32 bit words in header - min of 5)
	// zero out upper 4 bits by anding byte with 00001111
	internetHeaderLength := uint8(header[0] & byte(15))
	// 8 bit type of service
	typeOfService := uint8(header[1])
	// 16 bit datagram length
	totalLength := binary.LittleEndian.Uint16(header[2:4])
	// 16 bit identification number
	identification := binary.LittleEndian.Uint16(header[4:6])
	// 3 bit flags {0|DF|MF} (0 - reserved, DF = dont fragment, MF = more fragments)
	// the way to interpret the value of 'flags' variable is as follows:
	// == 0: no flags on
	// == 1: more fragments on
	// == 2: dont fragment is on
	// anything else is invalid, these two are mutually exclusive
	flags := uint8(header[6] >> 5)
	// 13 bit fragmentation offset
	// zero out upper 3 bits of first byte by anding it with 00011111
	fragOffset := binary.LittleEndian.Uint16([]byte{header[6] & byte(31), header[7]})
	// 8 bit time-to-live
	ttl := uint8(header[8])
	// 8 bit protocol value
	protocol := uint8(header[9])
	// 16 bit header checksum
	headerChecksum := binary.LittleEndian.Uint16(header[10:12])
	// 32 bit source ip
	sourceAddr := fmt.Sprintf("%d.%d.%d.%d", header[12], header[13], header[14], header[15])
	// 32 bit destination ip
	dstAddr := fmt.Sprintf("%d.%d.%d.%d", header[16], header[17], header[18], header[19])

	// TODO: HANDLE OPTIONS

	fmt.Printf("Version: %d\n", version)
	fmt.Printf("Internet Header Length: %d\n", internetHeaderLength)
	fmt.Printf("Type of Service: %d\n", typeOfService)
	fmt.Printf("Total Length: %d\n", totalLength)
	fmt.Printf("Identification: %d\n", identification)
	fmt.Printf("Flags: %d\n", flags)
	fmt.Printf("Fragment Offset: %d\n", fragOffset)
	fmt.Printf("TTL: %d\n", ttl)
	fmt.Printf("Protocol: %d\n", protocol)
	fmt.Printf("Header Checksum: %d\n", headerChecksum)
	fmt.Printf("Source IP: %s\n", sourceAddr)
	fmt.Printf("Destination IP: %s\n", dstAddr)
	return nil
}

func printRDTPHeader(p *rdtp.Packet) {
	fmt.Printf("Source Port: %d\n", p.SrcPort)
	fmt.Printf("Destination Port: %d\n", p.DstPort)
	fmt.Printf("Length: %d\n", p.Length)
	fmt.Printf("Checksum: %d\n", p.Checksum)
}
