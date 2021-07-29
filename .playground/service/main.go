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
	if err = svc.Run(); err != nil {
		log.Fatal(err)
	}
}
