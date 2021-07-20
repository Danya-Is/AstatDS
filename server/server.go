package main

import (
	"AstatDS"
	"bufio"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	state           = new(State)
	clientPortFlag  = flag.String("cp", ":8080", "flag for client communication")
	myPortFlag      = flag.String("p", "8081", "flag for technical communication")
	discoveryIpFlag = flag.String("d", "", "port belonging to one of already launched services in the cluster")
	clusterNameFlag = flag.String("c", "DefaultCluster", "name of the cluster to which service belongs")
	nodeNameFlag    = flag.String("n", "DefaultName", "name of the service")
	statePathFlag   = flag.String("s", "state", "state path")
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
			c.JSON(200, value)
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
	req := new(AstatDS.Request)
	err = json.Unmarshal(value, &req)
	state.KV[req.Key] = req.Value
	c.String(200, "OK")
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

	//TODO проверить чтобы пользоватеь указал все обязательные флаги и КОРРЕКТНО, вроде myPort, myClientPort

	if len(state.DiscoveryIpPort) > 0 {
		state.DiscoveryNodes()
	}

}

func WriteToDisk() {
	jsonstate, err := json.Marshal(state)
	if err != nil {
		log.Fatal(err)
	}
	stateEnc := base64.StdEncoding.EncodeToString(jsonstate)
	// overwriting content

	file, err := os.Create(state.StatePath)
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.WriteString(stateEnc)
	if err != nil {
		log.Fatal(err)
	}
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
}

func Loop() {
	for {
		state.CheckIps()
		state.CheckKV()

		str, _ := json.Marshal(state)
		if state.Hash != MD5(str) {
			WriteToDisk()
		}
	}
}

func listenNodes() {
	//TODO address
	ln, err := net.Listen("tcp", "127.0.0.1:"+state.MyPort)
	if err != nil {
		panic(err)
	}
	conn, err := ln.Accept()
	if err != nil {
		panic(err)
	}

	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		request := new(AstatDS.Request)
		json.Unmarshal([]byte(message), &request)

		switch request.Type {
		case AstatDS.GET_IPS:
			if _, ok := state.Ips[request.IP]; !ok {
				state.Ips[request.IP] = Node{
					time:   time.Now().String(),
					status: ACTIVATED,
				}
			}
			response, _ := json.Marshal(state.Ips)
			conn.Write([]byte(string(response) + "\n"))
		case AstatDS.GET_KV:
			response, _ := json.Marshal(state.Ips)
			conn.Write([]byte(string(response) + "\n"))
		case AstatDS.GET_IPS_HASH:
			str, _ := json.Marshal(state.Ips)
			response := MD5(str)
			conn.Write([]byte(response + "\n"))
		case AstatDS.GET_KV_HASH:
			str, _ := json.Marshal(state.KV)
			response := MD5(str)
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
	//TODO address
	sClient := &http.Server{
		Addr:           "0.0.0.0:" + state.MyClientPort,
		Handler:        clientRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	sClient.ListenAndServe()

	listenNodes()
}
