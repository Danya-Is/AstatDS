package commands

import (
	"AstatDS/client"
	"github.com/urfave/cli"
	"strings"
)

type response struct {
	value string
}

func NewGetCommand() *cli.Command {
	return &cli.Command{
		Name:      "get",
		Usage:     "get value from the store",
		UsageText: "get -key <key>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "key",
				Value:    "",
				Required: true,
			},
		},
		Action: get,
	}
}

func get(c *cli.Context) error {
	config, _ := client.ReadFromDisk()
	clientApi := client.New(config)

	respBody := clientApi.Get(strings.Trim(c.String("key"), "\n"))
	if respBody == nil {
		return client.KeyNotFoundError
	}
	print("responseBody: " + string(respBody) + "\n")
	print("\n")
	return nil
}
