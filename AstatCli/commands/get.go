package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"AstatDS/user"
	"github.com/urfave/cli/v2"
)

type getRequest struct {
	key string
}

type response struct {
	value string
}

func NewGetCommand() *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "get value from the store",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "key",
				Value: "",
			},
		},
		Action: get,
	}
}

func get(c *cli.Context) error {
	config, _ := user.Read()
	client := &http.Client{}
	for i := 0; i < len(config.Endpoints); i++ {
		reqBody, _ := json.Marshal(
			pushRequest{
				key: c.String("key"),
			})
		req, err := http.NewRequest("GET", config.Endpoints[i], bytes.NewReader(reqBody))
		if err != nil {
			//???
			continue
		}
		httpResp, err := client.Do(req)
		if err != nil {
			return err
		}
		respBody, _ := ioutil.ReadAll(httpResp.Body)
		resp := new(response)
		err = json.Unmarshal(respBody, resp)
		if err != nil {
			return err
		}
		fmt.Print(resp)
	}
	return nil
}
