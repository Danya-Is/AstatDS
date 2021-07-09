package commands

import (
	"io/ioutil"
	"os/user"

	"github.com/urfave/cli/v2"
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
	endpoints := c.String("endpoints")
	configData := []byte("Cluster name: " + name + "\nEndpoints: " + endpoints)
	usr, err := user.Current()
	if err != nil {
		return err
	}
	homeDir := usr.HomeDir
	ioutil.WriteFile(homeDir + "/AstatConfig", configData, 0777)
	return nil
}