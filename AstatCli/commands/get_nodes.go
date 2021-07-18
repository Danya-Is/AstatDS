package commands

import (
	"AstatDS/client"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
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
	resp := new(interface{})
	err := json.Unmarshal(respBody, resp)
	if err != nil {
		return err
	}
	fmt.Print(resp)
	return nil
}
