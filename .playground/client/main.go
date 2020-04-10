package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/adrianosela/rdtp"
)

func main() {
	c, err := net.Dial("unix", rdtp.DefaultRDTPServiceAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	c.Write([]byte("8.8.8.8"))

	fmt.Println("Anything written here will be sent to 8.8.8.8 over rdtp:")
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		if _, err := c.Write([]byte(text)[:len(text)-1]); err != nil {
			log.Fatal(err)
		}
	}
}
