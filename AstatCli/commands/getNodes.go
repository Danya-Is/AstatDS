package commands

import (
	"AstatDS/client"
	"fmt"
	"github.com/urfave/cli"
)

func NewGetNodesCommand() *cli.Command {
	return &cli.Command{
		Name:  "get-nodes",
		Usage: "get all nodes of the cluster",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Value: "Cluster",
				Usage: "name of the cluster",
			},
		},
		Action: getNodes,
	}
}

func getNodes(c *cli.Context) error {
	config, _ := client.ReadFromDisk()
	clientApi := client.New(config)
	respBody := clientApi.GetNodes()

	fmt.Print(string(respBody))
	return nil
}
