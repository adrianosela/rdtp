package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/adrianosela/rdtp"
	"github.com/pkg/errors"
)

func main() {
	addr := "10.0.0.94:22"

	c, err := rdtp.Dial(addr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Anything written here will be sent to %s over rdtp:\n", addr)
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		if _, err := c.Write([]byte(text)[:len(text)-1]); err != nil {
			log.Fatal(errors.Wrap(err, "Failed to write message"))
		}
	}
}
