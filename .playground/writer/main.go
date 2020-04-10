package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/adrianosela/rdtp"
)

func main() {
	addr := "192.168.1.77:29898"

	c, err := rdtp.Dial(addr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Anything written here will be sent to %s over rdtp:\n", addr)
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		if _, err := c.Write([]byte(text)[:len(text)-1]); err != nil {
			log.Fatal(err)
		}
	}
}
