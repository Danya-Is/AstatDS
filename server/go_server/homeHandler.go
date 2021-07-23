package go_server

import (
	"AstatDS"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"time"
)

func HomeGetHandler(c *gin.Context) {
	body := c.Request.Body
	data, err := ioutil.ReadAll(body)
	if err != nil {
		log.Fatal(err)
	}
	request := new(AstatDS.Request)
	err = json.Unmarshal(data, &request)
	switch request.Type {
	case AstatDS.GET_VALUE:
		key := request.Key
		value, ok := state.KV[key]
		fmt.Println(state)
		if ok {
			c.String(200, value.Value)
		} else {
			c.JSON(200, gin.H{"key": key, "value": "no value"})
		}
	case AstatDS.GET_NODES:
		data, _ := json.Marshal(state.Ips)
		c.String(200, string(data))

	}
}

func HomePostHandler(c *gin.Context) {
	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	if err != nil {
		log.Fatal(err)
	}
	req := new(AstatDS.Request)
	err = json.Unmarshal(value, &req)
	state.KV[req.Key] = Value{
		Time:  time.Now().Format(time_format),
		Value: req.Value,
	}
	fmt.Println(state)
	c.String(200, "OK")
}
