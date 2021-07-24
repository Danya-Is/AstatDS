package main

import (
	"log"
	"os"

	"AstatDS/AstatCli/commands"
	"github.com/urfave/cli"
)

func Start() error {
	app := cli.NewApp()
	app.Name = "AstatCli"
	app.Usage = "console application for Astat"
	app.Commands = []*cli.Command{
		commands.NewSetConfigCommand(),
		commands.NewPutCommand(),
		commands.NewGetCommand(),
		commands.NewGetNodesCommand(),
		commands.NewLSCommand(),
	}
	return app.Run(os.Args)
}

func main() {
	err := Start()
	if err != nil {
		log.Println(err)
	}
}
