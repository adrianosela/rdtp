package main

import (
	"github.com/adrianosela/rdtp/service"
	"github.com/pkg/errors"
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
		{
			Name:   "status",
			Usage:  "check the status of the rdtp service",
			Action: serviceStatusHandler,
		},
	},
}

func serviceStartHandler(ctx *cli.Context) error {
	svc, err := service.NewService()
	if err != nil {
		return errors.Wrap(err, "could not get rdtp service")
	}
	if err = svc.Start(); err != nil {
		return errors.Wrap(err, "could not start rdtp service")
	}
	return nil
}

func serviceStopHandler(ctx *cli.Context) error {
	// TODO
	return nil
}

func serviceStatusHandler(ctx *cli.Context) error {
	// TODO
	return nil
}
