package main

import (
	"github.com/adrianosela/rdtp/controller"
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
	},
}

func serviceStartHandler(ctx *cli.Context) error {
	c := controller.NewController()
	defer c.Shutdown()
	return c.Start()
}
