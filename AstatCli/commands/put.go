package commands

import (
	"AstatDS/client"
	"fmt"
	"github.com/urfave/cli"
	"log"
)

func NewPutCommand() *cli.Command {
	return &cli.Command{
		Name:      "put",
		Usage:     "push kv to the store",
		UsageText: "",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "key",
				Value:    "",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "value",
				Value:    "",
				Required: true,
			},
		},
		Action: put,
	}
}

func put(c *cli.Context) error {
	config, err := client.ReadFromDisk()
	if err != nil {
		log.Panicln(err)
		return err
	}
	clientApi := client.New(config)
	resp := clientApi.Put(c.String("key"), c.String("value"))
	fmt.Println(string(resp))
	return nil
}
