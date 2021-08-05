package commands

import (
	"errors"
	"log"
	"strings"

	"AstatDS/client"
	"github.com/urfave/cli"
)

func NewSetConfigCommand() *cli.Command {
	return &cli.Command{
		Name:  "set-config",
		Usage: "save cluster name and endpoints to file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Value: "NewCluster",
				Usage: "name of the cluster",
			},
			&cli.StringFlag{
				Name:  "endpoints",
				Value: "",
				Usage: "list of the endpoints\nExample: -endpoints 127.0.0.1:9080,127.0.0.1:9081",
			},
		},
		Action: setConfig,
	}
}

func setConfig(c *cli.Context) error {
	name := c.String("name")
	//TODO discover nodes
	endpoints := strings.Split(c.String("endpoints"), ",")
	if len(endpoints[len(endpoints)-1]) == 0 {
		log.Fatal(errors.New("Incorrect format\nDelete comma in the end of the list and do not use whitespaces"))
	}
	return (&client.Config{Cluster: name, Endpoints: endpoints}).Write()
}
