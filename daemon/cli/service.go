package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/adrianosela/rdtp/daemon"
	cli "gopkg.in/urfave/cli.v1"
)

var serviceCmds = cli.Command{
	Name:    "service",
	Aliases: []string{"s"},
	Usage:   "Manage rdtp service settings",
	Subcommands: []cli.Command{
		{
			Name:   "start",
			Usage:  "start the rdtp service",
			Action: serviceStartHandler,
		},
		{
			Name:   "stop",
			Usage:  "stop the rdtp service",
			Action: serviceStopHandler,
		},
	},
}

func serviceStartHandler(ctx *cli.Context) error {
	go runService()
	return nil
}

func serviceStopHandler(ctx *cli.Context) error {
	// TODO
	return nil
}

func runService() {
	c := daemon.NewController()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // block here until either SIGINT or SIGTERM is received

	c.Shutdown()
}
