package main

import (
	"AstatDS/server"
	"bufio"
	"encoding/json"
	"github.com/emirpasic/gods/maps/treemap"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type State struct {
	ClusterName  string `json:"clusterName"`
	MyIP         string `json:"myIP"`
	MyClientPort string `json:"myClientPort"`
	MyPort       string `json:"myPort"`
	DiscoveryIp  string `json:"discoveryIp"`
	NodeName     string `json:"nodeName"`
	StatePath    string `json:"statePath"`
}

var (
	KV  *treemap.Map
	Ips *treemap.Map
)

var mapMutex = sync.RWMutex{}
var StateHash string

const (
	ACTIVATED  = "activated"
	DEPRECATED = "deprecated"

	time_format = "2006-01-02 15:04:05 MST"
)

func (state *State) DiscoveryNodes() {
	conn, err := net.Dial("tcp", state.DiscoveryIp)
	if err != nil {
		log.Fatal(err)
	}
	str, err := json.Marshal(server.Request{
		Type: server.GET_IPS,
		IP:   state.MyIP + ":" + state.MyPort,
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Write([]byte(string(str) + "\n"))
	if err != nil {
		log.Fatal(err)
	}
	response, _ := bufio.NewReader(conn).ReadString('\n')
	err = Ips.FromJSON([]byte(response))
	if err != nil {
		log.Println(err)
	}

	err = conn.Close()
	if err != nil {
		log.Println(err)
	}
}

func UpdateNodeStatus(addr string, status string) {
	mapMutex.Lock()
	Ips.Put(addr, server.Node{
		Status: status,
		Time:   time.Now().Format(time_format),
	})
	mapMutex.Unlock()
}

func (state *State) CheckIps() {
	var ips []map[string]server.Node
	for addr, conn := range connections {
		mapMutex.Lock()
		intr, _ := Ips.Get(addr)
		mapMutex.Unlock()
		node := server.ConvertToNode(intr)
		status := node.Status

		if status == ACTIVATED && addr != state.MyIP+":"+state.MyPort {
			err := conn.SentRequest(server.GET_IPS_HASH, addr)
			if err != nil {
				log.Println(err)
				continue
			}
			response, _ := bufio.NewReader(conn.c).ReadString('\n')

			str, err := Ips.ToJSON()
			if err != nil {
				log.Println(err)
				continue
			}
			if strings.Compare(strings.Trim(response, "\n"), MD5(str)) != 0 {
				err = conn.SentRequest(server.GET_IPS, addr)
				if err != nil {
					log.Println(err)
					continue
				}
				response, _ := bufio.NewReader(conn.c).ReadString('\n')

				ip := new(map[string]server.Node)
				err = json.Unmarshal([]byte(response), &ip)
				if err != nil {
					log.Println(err)
				}
				ips = append(ips, *ip)
			}
		}
	}
	UpdateIps(ips)
}

func UpdateIps(ips []map[string]server.Node) {

	for _, m := range ips {
		for addr, node := range m {
			mapMutex.Lock()
			i, ok := Ips.Get(addr)
			if !ok {
				Ips.Put(addr, node)
			} else if n := server.ConvertToNode(i); n.Status != node.Status {
				t, err := time.Parse(time_format, node.Time)
				curT, err := time.Parse(time_format, n.Time)
				if err != nil {
					panic(err)
				}
				if t.After(curT) {
					Ips.Put(addr, node)
				}
			}
			mapMutex.Unlock()
		}
	}
}

func (state *State) CheckKV() {
	var kvs []map[string]server.Value
	for addr, conn := range connections {
		mapMutex.Lock()
		intr, _ := Ips.Get(addr)
		mapMutex.Unlock()
		node := server.ConvertToNode(intr)
		status := node.Status

		if status == ACTIVATED && addr != state.MyIP+":"+state.MyPort {
			err := conn.SentRequest(server.GET_KV_HASH, addr)
			if err != nil {
				log.Println("send req err")
				continue
			}
			response, _ := bufio.NewReader(conn.c).ReadString('\n')

			str, _ := KV.ToJSON()
			if strings.Compare(strings.Trim(response, "\n"), MD5(str)) != 0 {
				err := conn.SentRequest(server.GET_KV, addr)
				if err != nil {
					log.Println(err)
					continue
				}
				response, _ := bufio.NewReader(conn.c).ReadString('\n')

				kv := new(map[string]server.Value)
				err = json.Unmarshal([]byte(response), &kv)
				if err != nil {
					log.Println(err)
				}
				kvs = append(kvs, *kv)
			}
		}
	}
	UpdateKV(kvs)
}

func UpdateKV(kvs []map[string]server.Value) {
	mapMutex.Lock()
	for _, kv := range kvs {
		for k, v := range kv {
			i, ok := KV.Get(k)
			if !ok {
				KV.Put(k, v)
			} else if value := server.ConvertToValue(i); value.Value != v.Value {
				t, err := time.Parse(time_format, v.Time)
				curT, err := time.Parse(time_format, value.Time)
				if err != nil {
					log.Println(err)
				}
				if t.After(curT) {
					KV.Put(k, v)
				}
			}
		}
	}
	mapMutex.Unlock()
}
