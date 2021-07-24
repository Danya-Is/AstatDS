package commands

import (
	"AstatDS/client"
	"fmt"
	"github.com/urfave/cli"
)

func NewLSCommand() *cli.Command {
	return &cli.Command{
		Name:   "ls",
		Usage:  "get all kvs from the store",
		Action: ls,
	}
}

func ls(c *cli.Context) error {
	config, _ := client.ReadFromDisk()
	clientApi := client.New(config)
	respBody := clientApi.GetKVs()

	fmt.Print(string(respBody))
	return nil
}
