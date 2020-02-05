package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/adrianosela/rdtp"
	"github.com/pkg/errors"
)

func main() {
	// get raw network socket (AF_INET = IPv4)
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not get raw network socket"))
	}

	// specify whether messages to be sent have an ip header or need one to be added to them:
	// 0 - default ip header, 1 - custom ip header (we write it)
	syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 0)

	// define destination address (on loopback for now)
	addr := syscall.SockaddrInet4{Addr: [4]byte{127, 0, 0, 1}}

	fmt.Println("Anything written here will be sent over RDTP packets:")
	reader := bufio.NewReader(os.Stdin)
	for {
		// read user input
		text, _ := reader.ReadString('\n')

		// wrap it in a packet
		p, err := rdtp.NewPacket(uint16(14), uint16(15), []byte(text)[:len(text)-1])
		if err != nil {
			log.Println(errors.Wrap(err, "could not build rdtp packet for sending"))
		}

		// TODO:
		// if we have include header on (i.e. = 1), we have to wrap the rdtp packet
		// in an IPv4 datagram. In this step we:
		// - write the protocol number for RDTP = 157 (0x9D) (Unassigned as per https://en.wikipedia.org/wiki/List_of_IP_protocol_numbers)
	
		// send data to network socket
		if err = syscall.Sendto(fd, p.Serialize(), 0, &addr); err != nil {
			log.Fatal(errors.Wrap(err, "could not send data to network socket"))
		}
	}
}
