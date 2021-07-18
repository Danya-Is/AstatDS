package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"AstatDS"
	"github.com/gin-gonic/gin"
)

var (
	state = new(State)
)

func HomeGetHandler(c *gin.Context) {
	//TODO переделать под несколько типов запросов
	key := c.Query("key")
	value, ok := state.KV[key]
	if ok {
		c.JSON(200, gin.H{
			"key":   key,
			"value": value})
	} else {
		c.JSON(200, gin.H{"key": key, "value": "no value"})
	}
}

func HomePostHandler(c *gin.Context) {
	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	fmt.Println(string(value))
	if err != nil {
		log.Fatal(err)
	}
	var m interface{}
	err = json.Unmarshal(value, &m)
	data := m.(map[string]interface{})
	for k, v := range data {
		for k1 := range state.KV {
			if k1 == k {
				//TODO later should think about it more
				fmt.Println("this key already exists")
				break
			}
		}
		state.KV[k] = v
	}
	c.JSON(200, data)
}

func GetHandler(c *gin.Context) {
	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	fmt.Println(string(value))
	if err != nil {
		log.Fatal(err)
	}
	request := new(AstatDS.Request)
	json.Unmarshal(value, &request)
	switch request.Type {
	case AstatDS.GET_IPS:
		c.JSON(http.StatusOK, state.Ips)
	case AstatDS.GET_KV:
		//TODO сделать предварительную сверку по хэшу
		c.JSON(http.StatusOK, state.KV)
	}

}

func Init() {
	//читаем с диска

	//если стэйт пуст - ничего не делаем
	//если нет - записываем base64 -> json -> struct в стэйт state := State {...}

	//читаем флаги в стэйт
	//TODO проверить чтобы пользоватеь указал все обязательные флаги, вроде myPort, myClientPort

	state.DiscoveryNodes()
}

func WriteToDisk() {
	//записываем стэйт в файл
}

func Loop() {
	for {
		state.CheckIps()
		state.CheckKV()

		if state.hash != MD5(state) {
			WriteToDisk()
		}
	}
}

func listenNodes() {
	ln, _ := net.Listen("tcp", state.myPort)
	conn, _ := ln.Accept()

	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		request := new(AstatDS.Request)
		json.Unmarshal([]byte(message), &request)

		if request.Type == AstatDS.GET_IPS {
			response, _ := json.Marshal(state.Ips)
			conn.Write([]byte(string(response) + "\n"))
		} else if request.Type == AstatDS.GET_KV {
			response, _ := json.Marshal(state.Ips)
			conn.Write([]byte(string(response) + "\n"))
		}
	}
}

func main() {

	Init()
	go Loop()

	clientRouter := gin.Default()
	clientRouter.GET("/", HomeGetHandler)
	clientRouter.PUT("/", HomePostHandler)
	// r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	sClient := &http.Server{
		Addr:           state.myClientPort,
		Handler:        clientRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	sClient.ListenAndServe()

	listenNodes()
}
