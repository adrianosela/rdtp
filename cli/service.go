package main

import (
	"fmt"

	"github.com/adrianosela/rdtp/service"

	"github.com/pkg/errors"
	"github.com/takama/daemon"
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
	srv, err := daemon.New(ctx.App.Name, ctx.App.Description)
	if err != nil {
		return errors.Wrap(err, "could not get daemon")
	}
	status, err := srv.Install()
	if err != nil {
		return errors.Wrap(err, "could not install daemon")
	}
	fmt.Println(status)
	return nil
}

func serviceRemoveHandler(ctx *cli.Context) error {
	srv, err := daemon.New(ctx.App.Name, ctx.App.Description)
	if err != nil {
		return errors.Wrap(err, "could not get daemon")
	}
	status, err := srv.Remove()
	if err != nil {
		return errors.Wrap(err, "could not remove daemon")
	}
	fmt.Println(status)
	return nil
}

func serviceStartHandler(ctx *cli.Context) error {
	srv, err := daemon.New(ctx.App.Name, ctx.App.Description)
	if err != nil {
		return errors.Wrap(err, "could not get daemon")
	}
	status, err := srv.Start()
	if err != nil {
		return errors.Wrap(err, "could not start daemon")
	}
	fmt.Println(status)
	return nil
}

func serviceRunHandler(ctx *cli.Context) error {
	// // TODO: remove this listener
	// if err := c.Listen(uint16(15)); err != nil {
	// 	return errors.Wrap(err, "could not open new listener")
	// }

	srv, err := daemon.New(ctx.App.Name, ctx.App.Description)
	if err != nil {
		return errors.Wrap(err, "could not get daemon")
	}
	status, err := srv.Run(service.NewService())
	if err != nil {
		return errors.Wrap(err, "could not stop daemon")
	}
	fmt.Println(status)
	return nil
}

func serviceStopHandler(ctx *cli.Context) error {
	srv, err := daemon.New(ctx.App.Name, ctx.App.Description)
	if err != nil {
		return errors.Wrap(err, "could not get daemon")
	}
	status, err := srv.Stop()
	if err != nil {
		return errors.Wrap(err, "could not stop daemon")
	}
	fmt.Println(status)
	return nil
}

func serviceStatusHandler(ctx *cli.Context) error {
	srv, err := daemon.New(ctx.App.Name, ctx.App.Description)
	if err != nil {
		return errors.Wrap(err, "could not get daemon")
	}
	status, err := srv.Status()
	if err != nil {
		return errors.Wrap(err, "could not get status for daemon")
	}
	fmt.Println(status)
	return nil
}
