package service

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/adrianosela/rdtp"
	"github.com/pkg/errors"
)

// Client is a client of a host's TCP controller
type Client struct {
	svcConn net.Conn
}

// NewClient returns a new RDTP service client
func NewClient() (*Client, error) {
	svcConn, err := net.Dial("unix", rdtp.DefaultRDTPServiceAddr)
	if err != nil {
		return nil, errors.Wrap(err, "could not acquire RDTP service connection")
	}

	startShutdownSignalListener(svcConn)

	return &Client{
		svcConn: svcConn,
	}, nil
}

func startShutdownSignalListener(conn net.Conn) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func(c chan os.Signal) {
		sig := <-c
		log.Printf("[rdtp] shutting down - received signal: %s", sig)
		conn.Close()
		os.Exit(0)
	}(sigChan)
}

// Close closes an rdtp service client connection
func (c *Client) Close() error {
	log.Printf("[rdtp] shutting down - closed by client")
	return c.svcConn.Close()
}
