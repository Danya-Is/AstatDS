package commands

import (
	"AstatDS/client"
	"github.com/urfave/cli"
)

func NewPutCommand() *cli.Command {
	return &cli.Command{
		Name:  "put",
		Usage: "push kv to the store",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "key",
				Value: "",
			},
			&cli.StringFlag{
				Name:  "value",
				Value: "",
			},
		},
		Action: put,
	}
}

func put(c *cli.Context) error {
	config, err := client.ReadFromDisk()
	print("readed config\n")
	if err != nil {
		panic(err)
	}
	clientApi := client.New(config)
	return clientApi.Put(c.String("key"), c.String("value"))
}
