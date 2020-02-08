package main

import (
	"fmt"
	"os"

	cli "gopkg.in/urfave/cli.v1"
)

var version string // injected at build-time

func main() {
	app := cli.NewApp()

	app.Version = version
	app.EnableBashCompletion = true
	app.Usage = "rdtp (Reliable Data Transport Protocol) daemon"
	app.Flags = []cli.Flag{}

	app.Commands = []cli.Command{
		serviceCmds,
	}

	app.CommandNotFound = func(c *cli.Context, command string) {
		c.App.Run([]string{"help"})
		fmt.Printf("\ncommand \"%s\" does not exist\n", command)
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
