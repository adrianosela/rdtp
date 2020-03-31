package main

import (
	"fmt"
	"log"
	"math/rand"
	"syscall"
	"time"

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
	// define destination address
	addr := syscall.SockaddrInet4{Addr: [4]byte{192, 168, 1, 75}}

	fmt.Println("Flooding target with RDTP packets...")
	for {
		p, err := packet.NewPacket(uint16(rand.Intn(65534)+1), uint16(2), nil)
		if err != nil {
			log.Println(errors.Wrap(err, "could not build rdtp packet for sending"))
			continue
		}
		// send data to network socket
		if err = syscall.Sendto(fd, p.Serialize(), 0, &addr); err != nil {
			log.Fatal(errors.Wrap(err, "could not send data to network socket"))
		}

		time.Sleep(time.Second / 1000 * 10)
	}
}
