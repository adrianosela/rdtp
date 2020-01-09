package main

import (
	"flag"
)

var serve = flag.Bool("s", false, "Toggle server mode ON")

func main() {
	flag.Parse()

	if *serve {
		// TODO
	} else {
		// TODO
	}
}
