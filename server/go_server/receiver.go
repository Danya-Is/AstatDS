package main

import (
	"AstatDS/server"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

func sentHash(conn net.Conn, reqName string) error {
	var (
		str []byte
		err error
	)
	if reqName == server.GET_IPS_HASH {
		mapMutex.Lock()
		str, err = Ips.ToJSON()
		mapMutex.Unlock()
	} else {
		str, err = KV.ToJSON()
	}

	response := MD5(str)
	if err != nil {
		return err
	}
	_, err = conn.Write([]byte(response + "\n"))
	if err != nil {
		return err
	}
	return nil
}

func handle(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadBytes('\n')
		log.Println(string(message))
		if err != nil {
			log.Println("server disconnected")
			return
		}
		request := new(server.Request)
		err = json.Unmarshal(message, &request)
		if err != nil {
			log.Println(err)
		}

		if request.Type == server.GET_IPS {
			mapMutex.Lock()
			i, ok := Ips.Get(request.IP)
			if node := server.ConvertToNode(i); !ok || node.Status == DEPRECATED {
				log.Println("make activated")
				Ips.Put(request.IP, server.Node{
					Time:   time.Now().Format(time_format),
					Status: ACTIVATED,
				})
			}
			response, err := Ips.ToJSON()
			if err != nil {
				log.Println(err)
				return
			}
			mapMutex.Unlock()
			_, err = conn.Write([]byte(string(response) + "\n"))
			if err != nil {
				log.Println(err)
			}
		} else if request.Type == server.GET_KV {
			response, err := KV.ToJSON()
			if err != nil {
				log.Println(err)
				return
			}
			_, err = conn.Write([]byte(string(response) + "\n"))
			if err != nil {
				log.Println(err)
			}
		} else {
			err = sentHash(conn, request.Type)
			if err != nil {
				log.Println(err)
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
