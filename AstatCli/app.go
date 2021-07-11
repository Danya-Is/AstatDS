package main

import (
	"os"

	"AstatDS/AstatCli/commands"
	"github.com/urfave/cli/v2"
)

func main() {

	app := cli.NewApp()
	app.Name = "AstatCli"
	app.Usage = "console application for Astat"
	app.Commands = []*cli.Command{
		commands.NewSetConfigCommand(),
		commands.NewPushCommand(),
	}
	app.Run(os.Args)
}
