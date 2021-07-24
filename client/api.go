package client

import (
	"AstatDS/server"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type ClientApi interface {
	Get(key string) []byte
	Put(key string, value string) []byte
	GetNodes() []byte
	GetKVs() []byte
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

func (c *Client) Do(method string, reqBody []byte) ([]byte, error) {
	client := &http.Client{}
	for i := 0; i < len(c.Endpoints); i++ {
		req, err := http.NewRequest(method, "http://"+c.Endpoints[i], bytes.NewReader(reqBody))
		if err != nil {
			log.Println(err)
			continue
		}
		httpResp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			continue
		}
		respBody, _ := ioutil.ReadAll(httpResp.Body)
		return respBody, nil
	}
	return nil, NodesDontAnswer
}

func (c *Client) Get(key string) []byte {
	reqBody, err := json.Marshal(server.Request{
		Type: server.GET_VALUE,
		Key:  key,
	})
	if err != nil {
		log.Println(err)
		return nil
	}

	resp, err := c.Do("GET", reqBody)
	if err != nil {
		log.Println(err)
		return nil
	}
	return resp
}

func (c *Client) Put(key string, value string) []byte {
	reqBody, err := json.Marshal(
		server.Request{
			Type:  server.PUT_VALUE,
			Key:   key,
			Value: value,
		})
	if err != nil {
		log.Println(err)
		return nil
	}

	resp, err := c.Do("PUT", reqBody)
	if err != nil {
		log.Println(err)
		return nil
	}
	return resp
}

func (c *Client) GetNodes() []byte {
	reqBody, err := json.Marshal(server.Request{
		Type: server.GET_NODES,
	})
	if err != nil {
		log.Println(err)
		return nil
	}

	respBody, err := c.Do("GET", reqBody)
	if err != nil {
		log.Println(err)
		return nil
	}

	ips := make(map[string]server.Node)
	err = json.Unmarshal(respBody, &ips)
	if err != nil {
		return nil
	}

	resp := ""
	for ip, node := range ips {
		resp += ip + " - " + node.Status + "\n"
	}

	return []byte(resp)
}

func (c *Client) GetKVs() []byte {
	reqBody, err := json.Marshal(server.Request{
		Type: server.GET_KV,
	})
	if err != nil {
		log.Println(err)
		return nil
	}

	respBody, err := c.Do("GET", reqBody)
	if err != nil {
		log.Println(err)
		return nil
	}

	kvs := make(map[string]server.Value)
	err = json.Unmarshal(respBody, &kvs)
	if err != nil {
		return nil
	}

	resp := ""
	for k, v := range kvs {
		resp += k + " - " + v.Value + "\n"
	}

	return []byte(resp)
}
