package main

import (
	cli "gopkg.in/urfave/cli.v1"
)

var serviceCmds = cli.Command{
	Name:    "service",
	Aliases: []string{"s"},
	Usage:   "Manage rdtp service settings",
	Subcommands: []cli.Command{
		{
			Name:   "install",
			Usage:  "install the rdtp service",
			Action: serviceInstallHandler,
		},
		{
			Name:   "remove",
			Usage:  "remove the rdtp service",
			Action: serviceRemoveHandler,
		},
		{
			Name:   "start",
			Usage:  "start the rdtp service",
			Action: serviceStartHandler,
		},
		{
			Name:   "run",
			Usage:  "run the rdtp service",
			Action: serviceRunHandler,
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

func serviceInstallHandler(ctx *cli.Context) error {
	// TODO
	return nil
}

func serviceRemoveHandler(ctx *cli.Context) error {
	// TODO
	return nil
}

func serviceStartHandler(ctx *cli.Context) error {
	// TODO
	return nil
}

func serviceRunHandler(ctx *cli.Context) error {
	// TODO
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
