package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/adrianosela/rdtp"
	"github.com/pkg/errors"
)

var (
	serve = flag.Bool("s", false, "Toggle server mode ON")
)

func main() {
	flag.Parse()

	addr, err := net.ResolveIPAddr("ip", fmt.Sprintf("127.0.0.1"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not resolve IP address"))
	}

	if *serve {
		conn, err := net.ListenIP("ip:ip", addr)
		if err != nil {
			log.Fatal(errors.Wrap(err, "could not listen for IP"))
		}

		ipHeaderLen := 20

		// allocate buffer for packet
		buf := make([]byte, rdtp.MaxPacketBytes+ipHeaderLen) // ip header is 20 bytes
		for {
			// ReadFrom on an IPConn handles stripping the IP header
			// To examine ip header values we can use Read() or ReadString()
			ipPayloadSize, _, _ := conn.ReadFrom(buf)

			p, err := rdtp.Deserialize([]byte(buf)[:ipPayloadSize])
			if err != nil {
				log.Println(errors.Wrap(err, "could not build received rdtp packet"))
			}

			fmt.Println(string(p.Payload))
		}

	} else {
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
			p, err := rdtp.NewPacket(uint16(8080), uint16(8081), []byte(text)[:len(text)-1])
			if err != nil {
				log.Println(errors.Wrap(err, "could not build rdtp packet for sending"))
			}

			// send it to the server
			conn.Write(p.Serialize())
		}
	}
}
