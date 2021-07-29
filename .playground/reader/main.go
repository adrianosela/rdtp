package main

import (
	"io"
	"log"
	"net"

	"github.com/adrianosela/rdtp"
)

func main() {
	addr := "10.0.0.94:22"

	l, err := rdtp.Listen(addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	log.Printf("Listening for new connections on rdtp address %s\n", addr)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			break
		}
		log.Printf("Accepted new connection from %s\n", conn.RemoteAddr())

		go func(c net.Conn) {
			for {
				buf := make([]byte, 1024)
				n, err := c.Read(buf)
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Printf("ERROR: %s\n", err)
				}
				log.Printf("[%s] %s", c.RemoteAddr(), string(buf[:n]))
			}
		}(conn)
	}
}
