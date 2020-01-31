package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/adrianosela/rdtp"
	"github.com/pkg/errors"
)

func main() {
	addr, err := net.ResolveIPAddr("ip", "192.168.1.71")
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not resolve IP address"))
	}

	conn, err := net.DialIP("ip:ip", nil, addr)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not dial IP"))
	}

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

		// send it to the server
		conn.Write(p.Serialize())
	}
}
