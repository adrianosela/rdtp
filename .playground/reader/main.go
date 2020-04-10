package main

import (
	"fmt"
	"log"

	"github.com/adrianosela/rdtp"
)

func main() {
	addr := "192.168.1.77"

	c, err := rdtp.Dial(addr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Anything written here will be sent to %s over rdtp:\n", addr)
	buf := make([]byte, 15000)
	for {
		n, err := c.Read(buf)
		if err != nil {
			log.Println(err)
		}

		fmt.Println(string(buf[:n]))
	}
}
