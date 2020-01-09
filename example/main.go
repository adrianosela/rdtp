package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"

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
		rdr := bufio.NewReader(conn)
		for {
			data, _, _ := rdr.ReadLine()
			fmt.Println(string([]byte(data)[20:]))
		}
	} else {
		conn, err := net.DialIP("ip:ip", nil, addr)
		if err != nil {
			log.Fatal(errors.Wrap(err, "could not dial IP"))
		}

		conn.Write([]byte("eat my shorts\n"))
	}
}

