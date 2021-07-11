package commands

import (
	"bytes"
	"encoding/json"
	"net/http"

	"AstatDS/user"
	"github.com/urfave/cli/v2"
)

type pushRequest struct {
	key   string
	value string
}

func NewPushCommand() *cli.Command {
	return &cli.Command{
		Name:  "push",
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
		Action: push,
	}
}

func push(c *cli.Context) error {
	config, _ := user.Read()
	client := &http.Client{}
	for i := 0; i < len(config.Endpoints); i++ {
		reqBody, _ := json.Marshal(
			pushRequest{
				key:   c.String("key"),
				value: c.String("value"),
			})
		req, err := http.NewRequest("PUSH", config.Endpoints[i], bytes.NewReader(reqBody))
		if err != nil {
			//???
			continue
		}
		_, err = client.Do(req)
		if err != nil {
			return err
		}
	}
	return nil
}
