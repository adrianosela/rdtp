package main

import (
	"log"

	"github.com/adrianosela/rdtp/service"
)

func main() {
	svc, err := service.NewService()
	if err != nil {
		log.Fatal(err)
	}

	err = svc.Start()
	if err != nil {
		log.Fatal(err)
	}
}
