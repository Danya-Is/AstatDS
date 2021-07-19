package client

import (
	"AstatDS"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type ClientApi interface {
	Get(key string) []byte
	Put(key string, value string) error
	GetNodes() []byte
}

type Client struct {
	Endpoints []string
	Cluster   string
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
		reqBody, _ := json.Marshal(AstatDS.Request{
			Type: AstatDS.GET_VALUE,
			Key:  key,
		})
		req, err := http.NewRequest("GET", "http://"+c.Endpoints[i], bytes.NewReader(reqBody))
		if err != nil {
			log.Fatal(err)
			continue
		}
		httpResp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
			continue
		}
		respBody, _ := ioutil.ReadAll(httpResp.Body)
		return respBody
	}
	return nil
}

func (c *Client) Put(key string, value string) error {
	client := &http.Client{}
	for i := 0; i < len(c.Endpoints); i++ {
		reqBody, err := json.Marshal(
			AstatDS.Request{
				Type:  AstatDS.PUT_VALUE,
				Key:   key,
				Value: value,
			})
		print("created reqBody\n")
		req, err := http.NewRequest("POST", "http://"+c.Endpoints[i], bytes.NewReader(reqBody))
		print("created req\n")
		if err != nil {
			log.Fatal(err)
			continue
		}
		_, err = client.Do(req)
		print("did req\n")
		if err != nil {
			log.Fatal(err)
			continue
		}
		return nil
	}
	return NodesDontAnswer
}

func (c *Client) GetNodes() []byte {
	client := &http.Client{}
	for i := 0; i < len(c.Endpoints); i++ {
		reqBody, _ := json.Marshal(AstatDS.Request{
			Type: AstatDS.GET_NODES,
		})
		req, err := http.NewRequest("GET", c.Endpoints[i], bytes.NewReader(reqBody))
		if err != nil {
			log.Fatal(err)
			continue
		}
		httpResp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
			continue
		}
		respBody, _ := ioutil.ReadAll(httpResp.Body)
		return respBody
	}
	return nil
}
