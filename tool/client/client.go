package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/adrianosela/rdtp"
	"github.com/adrianosela/rdtp/packet"
	"github.com/pkg/errors"
)

func main() {
	// get raw network socket (AF_INET = IPv4)
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, rdtp.IPProtoRDTP)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not get raw network socket"))
	}

	// define destination address (on loopback for now)
	addr := syscall.SockaddrInet4{Addr: [4]byte{192, 168, 1, 75}}

	fmt.Println("Anything written here will be sent over RDTP packets:")
	reader := bufio.NewReader(os.Stdin)
	for {
		// read user input
		text, _ := reader.ReadString('\n')

		// wrap it in a packet
		p, err := packet.NewPacket(uint16(14), uint16(2), []byte(text)[:len(text)-1])
		if err != nil {
			log.Println(errors.Wrap(err, "could not build rdtp packet for sending"))
		}

		// send data to network socket
		if err = syscall.Sendto(fd, p.Serialize(), 0, &addr); err != nil {
			log.Fatal(errors.Wrap(err, "could not send data to network socket"))
		}
	}
}
