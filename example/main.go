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

		buf := make([]byte, rdtp.MaxPacketBytes+20) // ip header is 20 bytes
		for {
			// ReadFrom on an IPConn handles stripping the IP header
			// To examine ip header values we can use Read() or ReadString()
			ipDatagramLength, _, _ := conn.ReadFrom(buf)
			fmt.Println(string([]byte(buf)[:ipDatagramLength]))
		}

	} else {
		conn, err := net.DialIP("ip:ip", nil, addr)
		if err != nil {
			log.Fatal(errors.Wrap(err, "could not dial IP"))
		}

		fmt.Println("Anything written here will be sent over IP packets:")
		reader := bufio.NewReader(os.Stdin)
		for {
			text, _ := reader.ReadString('\n')
			conn.Write([]byte(text)[:len(text)-1])
		}
	}
}
