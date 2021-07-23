package go_server

import (
	"AstatDS"
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
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

func ReadFlags() {
	state.MyClientPort = *clientPortFlag
	state.MyPort = *myPortFlag
	state.DiscoveryIp = *discoveryIpFlag
	state.ClusterName = *clusterNameFlag
	state.NodeName = *nodeNameFlag
	state.StatePath = *statePathFlag
	state.MyIP = *ipFlag
}

func Init() {

	file, err := ReadFile(*statePathFlag)
	if err != nil {
		log.Fatal(err)
	}

	if len(file) > 0 {
		ReadState(file)
	} else {
		state.KV = make(map[string]Value)
		state.Ips = make(map[string]Node)
	}

	ReadFlags()
	checkFlags()

	UpdateNodeStatus(state.MyIP+":"+state.MyPort, ACTIVATED)

	if len(state.DiscoveryIp) > 0 {
		state.DiscoveryNodes()
	}
}

func Loop() {
	for {
		UpdateConnections()
		state.CheckIps()
		state.CheckKV()

		mapMutex.Lock()
		str, _ := json.Marshal(state)
		mapMutex.Unlock()

		if StateHash != MD5(str) {
			StateHash = MD5(str)
			WriteToDisk()
		}

		time.Sleep(1000)
	}
}

func handle(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			fmt.Println("server disconnected")
			return
		}
		request := new(AstatDS.Request)
		err = json.Unmarshal(message, &request)
		if err != nil {
			log.Println(err)
		}

		switch request.Type {
		case AstatDS.GET_IPS:
			mapMutex.Lock()
			if _, ok := state.Ips[request.IP]; !ok {
				log.Println("add IP " + request.IP)
				state.Ips[request.IP] = Node{
					Time:   time.Now().Format(time_format),
					Status: ACTIVATED,
				}
			} else if state.Ips[request.IP].Status == DEPRECATED {
				//log.Println("change IP " + request.IP)
				state.Ips[request.IP] = Node{
					Time:   time.Now().Format(time_format),
					Status: ACTIVATED,
				}
			}
			response, err := json.Marshal(state.Ips)
			if err != nil {
				log.Println(err)
				return
			}
			mapMutex.Unlock()
			_, err = conn.Write([]byte(string(response) + "\n"))
			if err != nil {
				log.Println(err)
			}
		case AstatDS.GET_KV:
			response, err := json.Marshal(state.KV)
			if err != nil {
				log.Println(err)
				return
			}
			_, err = conn.Write([]byte(string(response) + "\n"))
			if err != nil {
				log.Println(err)
			}
		case AstatDS.GET_IPS_HASH:
			mapMutex.Lock()
			str, err := json.Marshal(state.Ips)
			mapMutex.Unlock()
			response := MD5(str)
			if err != nil {
				log.Println(err)
				return
			}
			_, err = conn.Write([]byte(response + "\n"))
			if err != nil {
				log.Println(err)
			}
		case AstatDS.GET_KV_HASH:
			str, err := json.Marshal(state.KV)
			if err != nil {
				log.Println(err)
				return
			}
			response := MD5(str)
			_, err = conn.Write([]byte(response + "\n"))
			if err != nil {
				log.Println(err)
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
