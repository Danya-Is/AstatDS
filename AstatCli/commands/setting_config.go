package commands

import (
	"strings"

	"AstatDS/client"
	"github.com/urfave/cli/v2"
)

func NewSetConfigCommand() *cli.Command {
	return &cli.Command{
		Name:  "set-client",
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
	endpoints := strings.Split(c.String("endpoints"), ",")
	return (&client.Config{Cluster: name, Endpoints: endpoints}).Write()
}
