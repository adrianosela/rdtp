package controller

import (
	"fmt"
	"log"
)

// Worker handles the local transport layer processing
// for a single process-process communication
type Worker struct {
	id     uint16
	rxChan chan []byte
}

// NewWorker returns an RDTP Worker struct
func NewWorker() (*Worker, error) {
	w := &Worker{
		id:     0, // reserved port number 0
		rxChan: make(chan []byte),
	}

	go w.reader()

	return w, nil
}

func (w *Worker) reader() error {
	for {
		select {
		case message, ok := <-w.rxChan:
			if !ok {
				return fmt.Errorf("worker receive channel closed")
			}
			// TODO
			log.Printf("[WORKER %d] received: %s", w.id, string(message))
		}
	}
}

// Kill shuts down a worker
func (w *Worker) Kill() error {
	close(w.rxChan)
	return nil
}
