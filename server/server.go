package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"log"
	"net"
	"net/http"
	"time"
	"flag"
	"AstatDS"
	"github.com/gin-gonic/gin"
	"encoding/base64"
)

var (
	state = new(State)
	clientPortFlag = flag.String("cp", ":8080", "flag for client communication")
	myPortFlag = flag.String("p", ":8081", "flag for technical communication")
	discoveryIpFlag = flag.String("d", ":8082", "port belonging to one of already launched services in the cluster")
	clusterNameFlag = flag.String("c", "DefaultCluster", "name of the cluster to which service belongs")
	nodeNameFlag = flag.String("n", "DefaultName", "name of the service")
	statePathFlag = flag.String("s", "state", "state path")
)

func HomeGetHandler(c *gin.Context) {
	body := c.Request.Body
	data, err := ioutil.ReadAll(body)
	fmt.Println(string(data))
	if err != nil {
		log.Fatal(err)
	}
	request := new(AstatDS.Request)
	err = json.Unmarshal(data, &request)
	switch request.Type {
	case AstatDS.GET_VALUE:
		key := request.Key
		value, ok := state.KV[key]
		if ok {
			fmt.Println(value)
			c.JSON(200, gin.H{
				"key":   key,
				"value": value})
		} else {
			c.JSON(200, gin.H{"key": key, "value": "no value"})
		}
	case AstatDS.GET_NODES:
		c.JSON(200, state.Ips)

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

func Init() {
	flag.Parse()
	if _, err := os.Stat(*statePathFlag); os.IsNotExist(err) {
		os.Create(*statePathFlag) // create file if it isn't exist
	}
	file, err := ioutil.ReadFile(*statePathFlag)
	if err != nil {
		log.Fatal(err)
	}
	if len(file) > 0 { // that means if file is not empty
		fDec, err := base64.StdEncoding.DecodeString(string(file))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(fDec)) 
		// jsDec, _ := json.Marshal(fDec) // idk if this really has to be here, we can assume that we alwsys have proper json, dont we?
		// fmt.Println(string(jsDec)) // fDec and jsDec are the same... except '\n' thing
		err = json.Unmarshal(fDec, &state)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		state.KV = make(map[string]interface{})
		state.Ips = make(map[string]interface{})
	}
	fmt.Println(state)
	state.MyClientPort = *clientPortFlag
	state.MyPort = *myPortFlag
	state.DiscoveryIpPort = *discoveryIpFlag
	state.ClusterName = *clusterNameFlag
	state.NodeName = *nodeNameFlag
	state.StatePath = *statePathFlag
	fmt.Println(state)
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

		if state.Hash != StateMD5(state) {
			WriteToDisk()
		}
	}
}

func listenNodes() {
	ln, _ := net.Listen("tcp", state.MyPort)
	conn, _ := ln.Accept()

	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		request := new(AstatDS.Request)
		json.Unmarshal([]byte(message), &request)

		switch request.Type {
		case AstatDS.GET_IPS:
			response, _ := json.Marshal(state.Ips)
			conn.Write([]byte(string(response) + "\n"))
		case AstatDS.GET_KV:
			response, _ := json.Marshal(state.Ips)
			conn.Write([]byte(string(response) + "\n"))
		case AstatDS.GET_IPS_HASH:
			response := MD5(state.Ips)
			conn.Write([]byte(response + "\n"))
		case AstatDS.GET_KV_HASH:
			response := MD5(state.KV)
			conn.Write([]byte(response + "\n"))
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
		Addr:           state.MyClientPort,
		Handler:        clientRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	sClient.ListenAndServe()

	listenNodes()
}
