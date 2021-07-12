package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type ClientApi interface {
	Get()
	Put()
}

type Client struct {
	Endpoints []string
	Cluster   string
}

type Request struct {
	key   string
	value []byte
}

func New(config *Config) Client {
	return Client{
		Endpoints: config.Endpoints,
		Cluster:   config.Cluster,
	}
}

func (c *Client) Get(key string) []byte {
	client := &http.Client{}
	for i := 0; i < len(c.Endpoints); i++ {
		req, err := http.NewRequest("GET", c.Endpoints[i], bytes.NewReader([]byte(key)))
		if err != nil {
			//???
			continue
		}
		httpResp, err := client.Do(req)
		if err != nil {
			//???
			continue
		}
		respBody, _ := ioutil.ReadAll(httpResp.Body)
		return respBody
	}
	return nil
}

func (c *Client) Put(key string, value []byte) error {
	client := &http.Client{}
	for i := 0; i < len(c.Endpoints); i++ {
		reqBody, _ := json.Marshal(
			Request{
				key:   key,
				value: value,
			})
		req, err := http.NewRequest("POST", c.Endpoints[i], bytes.NewReader(reqBody))
		if err != nil {
			//???
			continue
		}
		_, err = client.Do(req)
		if err != nil {
			//???
			continue
		}
		return nil
	}
	return KeyNotFoundError
}
