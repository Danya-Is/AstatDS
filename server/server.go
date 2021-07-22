package main

import (
	"AstatDS"
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
	"strconv"
	"strings"
	"time"
)

var (
	state           = new(State)
	clientPortFlag  = flag.String("cp", "", "flag for client communication")
	myPortFlag      = flag.String("p", "", "flag for technical communication")
	discoveryIpFlag = flag.String("d", "", "ip belonging to one of already launched services in the cluster")
	ipFlag          = flag.String("i", "", "my ip")
	clusterNameFlag = flag.String("c", "DefaultCluster", "name of the cluster to which service belongs")
	nodeNameFlag    = flag.String("n", "", "name of the service")
	statePathFlag   = flag.String("s", "", "state path")

	connections map[string]net.Conn
)

func checkFlags() {
	//ClientPort and myPort
	if *clientPortFlag == "" {
		fmt.Println("Error: clientPortFlag isn't specified")
		os.Exit(2)
	}
	num, err := strconv.Atoi(*clientPortFlag)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	if (num != 80 && num != 81) && (num < 1024 || num > 49151) {
		fmt.Println("Error: clientPortFlag isn't correct")
		os.Exit(2)
	}
	if *myPortFlag == "" {
		fmt.Println("Error: myPortFlag isn't specified")
		os.Exit(2)
	}
	num, err = strconv.Atoi(*myPortFlag)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	if (num != 80 && num != 81) && (num < 1024 || num > 49151) {
		fmt.Println("Error: myPortFlag isn't correct")
		os.Exit(2)
	}
	// my ip
	if *ipFlag == "" {
		fmt.Println("Error: ipFlag isn't specified")
		os.Exit(2)
	}
	num4 := strings.Split(*ipFlag, ".")
	if len(num4) != 4 {
		fmt.Println("Error: incorrect ipFlag")
		os.Exit(2)
	}
	for i := range num4 {
		num, err = strconv.Atoi(num4[i])
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		if num < 0 || num > 255 {
			fmt.Println("Error: incorrect ipFlag")
			os.Exit(2)
		}
	}
	/* discoveryIpFlag (TODO more thorough check)
	if *discoveryIpFlag == "" {
		fmt.Println("Error: discoveryIpFlag isn't specified")
		os.Exit(2)
	}
	num4 = strings.Split(*ipFlag, ".")
	if len(num4) != 4 {
		fmt.Println("Error: incorrect discoveryIpFlag")
	}
	for i:= 0; i < 3; i++ {
		num, err = strconv.Atoi(num4[i])
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		if num < 0 || num > 255 {
			fmt.Println("Error: incorrect ipFlag")
			os.Exit(2)
		}
	}
	num2 := strings.Split(*ipFlag, ":")
	for i := range num2 {
		num, err = strconv.Atoi(num4[i])
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		if (num < 0 || num > 255) && (i == 0) {
			fmt.Println("Error: incorrect discoveryIpFlag")
			os.Exit(2)
		}
		if (num != 80 && num != 81) && (num < 1024 || num > 49151) && (i == 1) {
			fmt.Println("Error: incorrect discoveryIpFlag")
			os.Exit(2)
		}
	}*/
	// statePathFlag
	if *statePathFlag == "" {
		fmt.Println("Error: statePath isn't specified")
		os.Exit(2)
	}
	// nodeNameFlag
	if *nodeNameFlag == "" {
		fmt.Println("Error: nodeName isn't specified")
		os.Exit(2)
	}
}

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
	state.KV[req.Key] = Value{
		time:  time.Now().Format(time_format),
		value: req.Value,
	}
	c.String(200, "OK")
}

func Init() {
	home, _ := os.UserHomeDir()
	flag.Parse()
	if _, err := os.Stat(home + "/" + *statePathFlag); os.IsNotExist(err) {
		os.Create(home + "/" + *statePathFlag) // create file if it isn't exist
	}
	fmt.Println(home)
	file, err := ioutil.ReadFile(home + "/" + *statePathFlag)
	if err != nil {
		log.Fatal(err)
	}
	if len(file) > 0 { // that means if file is not empty
		data := strings.Split(string(file), "\n")
		hashDec, err := base64.StdEncoding.DecodeString(data[0])
		if err != nil {
			log.Fatal(err)
		}
		StateHash = string(hashDec)
		fmt.Println("StateHash: " + StateHash)
		stateDec, err := base64.StdEncoding.DecodeString(data[1])
		fmt.Println("State: " + string(stateDec))
		err = json.Unmarshal(stateDec, &state)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		state.KV = make(map[string]Value)
		state.Ips = make(map[string]Node)
	}
	state.MyClientPort = *clientPortFlag
	state.MyPort = *myPortFlag
	state.DiscoveryIp = *discoveryIpFlag
	state.ClusterName = *clusterNameFlag
	state.NodeName = *nodeNameFlag
	state.StatePath = *statePathFlag
	state.MyIP = *ipFlag
	fmt.Println(state)
	checkFlags()

	if len(state.DiscoveryIp) > 0 {
		state.DiscoveryNodes()
	}
	Connections()

}

func Connections() {
	for addr := range state.Ips {
		var err error
		connections[addr], err = net.Dial("tcp", addr)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func WriteToDisk() {
	home, _ := os.UserHomeDir()
	jsonstate, err := json.Marshal(state)
	if err != nil {
		log.Fatal(err)
	}
	stateEnc := base64.StdEncoding.EncodeToString(jsonstate)
	// overwriting content

	file, err := os.Create(home + "/" + state.StatePath)
	if err != nil {
		log.Fatal(err)
	}
	StateHashEnc := base64.StdEncoding.EncodeToString([]byte(StateHash))
	_, err = file.WriteString(StateHashEnc + "\n") // write StateHash as a first string
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.WriteString(stateEnc) // write State encoded in base64 as a second string
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
		if StateHash != MD5(str) {
			StateHash = MD5(str)
			WriteToDisk()
		}
	}
}

func handle(conn net.Conn) {
	for {
		//message, err := bufio.NewReader(conn).ReadBytes('\n')
		var message []byte
		_, err := conn.Read(message)
		if err != nil {
			fmt.Println("server disconnected")
			return
		}
		fmt.Println("message" + string(message) + ".")
		request := new(AstatDS.Request)
		err = json.Unmarshal([]byte(message), &request)
		if err != nil {
			return
		}

		fmt.Println(request)

		switch request.Type {
		case AstatDS.GET_IPS:
			if _, ok := state.Ips[request.IP]; !ok {
				state.Ips[request.IP] = Node{
					time:   time.Now().Format(time_format),
					status: ACTIVATED,
				}
			}
			response, _ := json.Marshal(state.Ips)
			_, err = conn.Write([]byte(string(response) + "\n"))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("got GET_IPS")
		case AstatDS.GET_KV:
			response, _ := json.Marshal(state.Ips)
			_, err = conn.Write([]byte(string(response) + "\n"))
			if err != nil {
				return
			}
		case AstatDS.GET_IPS_HASH:
			str, _ := json.Marshal(state.Ips)
			response := MD5(str)
			_, err = conn.Write([]byte(response + "\n"))
			if err != nil {
				return
			}
		case AstatDS.GET_KV_HASH:
			str, _ := json.Marshal(state.KV)
			response := MD5(str)
			_, err = conn.Write([]byte(response + "\n"))
			if err != nil {
				return
			}
		}
	}
}

func listenNodes(c chan int) {
	ln, err := net.Listen("tcp", state.MyIP+":"+state.MyPort)
	if err != nil {
		panic(err)
	}
	defer func(ln net.Listener) {
		err := ln.Close()
		c <- 0
		if err != nil {
			log.Println(err)
		}
	}(ln)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("new conn")
		go handle(conn)
	}
}

func main() {

	c := make(chan int)
	Init()
	go Loop()
	go listenNodes(c)

	clientRouter := gin.Default()
	clientRouter.GET("/", HomeGetHandler)
	clientRouter.PUT("/", HomePostHandler)
	sClient := &http.Server{
		Addr:           state.MyIP + ":" + state.MyClientPort,
		Handler:        clientRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := sClient.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

	res := <-c
	os.Exit(res)
}
