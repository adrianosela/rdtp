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

	// specify whether send messages will have an ip header
	// or need one to be added to them:
	// 1 - custom ip header
	// 0 - default ip hedaer
	syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 0)

	// define destination address
	addr := syscall.SockaddrInet4{Addr: [4]byte{128, 189, 200, 255}}

	fmt.Println("Anything written here will be sent over IP packets:")
	reader := bufio.NewReader(os.Stdin)
	for {
		// read user input
		text, _ := reader.ReadString('\n')

		// wrap it in a packet
		p, err := rdtp.NewPacket(uint16(14), uint16(15), []byte(text)[:len(text)-1])
		if err != nil {
			log.Println(errors.Wrap(err, "could not build rdtp packet for sending"))
		}

		// TODO: wrap rdtp packet in IPv4 datagram and write protocol number to IP datagram

		// send it to socket
		if err = syscall.Sendto(fd, p.Serialize(), 0, &addr); err != nil {
			log.Fatal(errors.Wrap(err, "could not send data to socket"))
		}
	}
}
