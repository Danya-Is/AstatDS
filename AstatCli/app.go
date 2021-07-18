package main

import (
	"os"

	"AstatDS/AstatCli/commands"
	"github.com/urfave/cli/v2"
)

func Start() error {
	app := cli.NewApp()
	app.Name = "AstatCli"
	app.Usage = "console application for Astat"
	app.Commands = []*cli.Command{
		commands.NewSetConfigCommand(),
		commands.NewPutCommand(),
		commands.NewGetCommand(),
	}
	return app.Run(os.Args)
}

func main() {
	err := Start()
	if err != nil {
		panic(err)
	}
}
