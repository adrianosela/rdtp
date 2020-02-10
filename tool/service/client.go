package main

import (
	"fmt"
	"log"

	"github.com/adrianosela/rdtp"
)

func main() {
	rdtpConn, err := rdtp.Dial([]byte{127, 0, 0, 1})
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1500)

	for {
		n, err := rdtpConn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(buf[:n]))
	}
}
