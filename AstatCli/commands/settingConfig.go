package commands

import (
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
				Usage: "list of the endpoints",
			},
		},
		Action: setConfig,
	}
}

func setConfig(c *cli.Context) error {
	name := c.String("name")
	//TODO discover nodes
	endpoints := strings.Split(c.String("endpoints"), ",")
	return (&client.Config{Cluster: name, Endpoints: endpoints}).Write()
}