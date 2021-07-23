package go_server

import (
	"AstatDS"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

type State struct {
	KV           map[string]Value `json:"kv"`
	Ips          map[string]Node  `json:"ips"`
	ClusterName  string           `json:"clusterName"`
	MyIP         string           `json:"myIP"`
	MyClientPort string           `json:"myClientPort"`
	MyPort       string           `json:"myPort"`
	DiscoveryIp  string           `json:"discoveryIp"`
	NodeName     string           `json:"nodeName"`
	StatePath    string           `json:"statePath"`
}

var StateHash string

type Node struct {
	Time   string `json:"time"`
	Status string `json:"status"`
}

type Value struct {
	Time  string `json:"time"`
	Value string `json:"value"`
}

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
	str, err := json.Marshal(AstatDS.Request{
		Type: AstatDS.GET_IPS,
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
	fmt.Println(response)
	err = json.Unmarshal([]byte(response), &state.Ips)
	if err != nil {
		log.Println(err)
	}
	log.Println("nodes discovered")
	fmt.Println(state)
	err = conn.Close()
	if err != nil {
		log.Println(err)
	}
}

func UpdateNodeStatus(addr string, status string) {
	mapMutex.Lock()
	state.Ips[addr] = Node{
		Status: status,
		Time:   time.Now().Format(time_format),
	}
	mapMutex.Unlock()
}

func (state *State) CheckIps() {
	var ips []map[string]Node
	for addr, conn := range connections {
		mapMutex.Lock()
		if state.Ips[addr].Status == ACTIVATED && addr != state.MyIP+":"+state.MyPort {
			mapMutex.Unlock()
			str, _ := json.Marshal(AstatDS.Request{
				Type: AstatDS.GET_IPS_HASH,
				IP:   state.MyIP + ":" + state.MyPort,
			})

			if conn == nil {
				UpdateNodeStatus(addr, DEPRECATED)
				continue
			}

			_, err := fmt.Fprintf(conn, string(str)+"\n")
			if err != nil {
				log.Println(err)
				UpdateNodeStatus(addr, DEPRECATED)
				continue
			}
			response, _ := bufio.NewReader(conn).ReadString('\n')

			if response != MD5(str) {
				str, _ := json.Marshal(AstatDS.Request{
					Type: AstatDS.GET_IPS,
					IP:   state.MyIP + ":" + state.MyPort,
				})

				_, err := fmt.Fprintf(conn, string(str)+"\n")
				if err != nil {
					log.Println(err)
					UpdateNodeStatus(addr, DEPRECATED)
					continue
				}
				response, _ := bufio.NewReader(conn).ReadString('\n')
				ip := new(map[string]Node)
				err = json.Unmarshal([]byte(response), &ip)
				if err != nil {
					log.Println(err)
				}
				ips = append(ips, *ip)
			}
		} else {
			mapMutex.Unlock()
		}

	}

	UpdateIps(ips)
}

func UpdateIps(ips []map[string]Node) {

	for _, m := range ips {
		for addr, node := range m {
			mapMutex.Lock()
			if _, ok := state.Ips[addr]; !ok {
				state.Ips[addr] = node
			} else if n, _ := state.Ips[addr]; n.Status != node.Status {
				t, err := time.Parse(time_format, node.Time)
				curT, err := time.Parse(time_format, n.Time)
				if err != nil {
					panic(err)
				}
				if t.After(curT) {
					state.Ips[addr] = node
				}
			}
			mapMutex.Unlock()
		}
	}
}

func (state *State) CheckKV() {
	var kvs []map[string]Value
	for addr, conn := range connections {
		mapMutex.Lock()
		log.Println("tuck-tuck to " + addr)
		//if state.Ips[addr].Status == ACTIVATED && addr != state.MyIP+":"+state.MyPort {
		if addr != state.MyIP+":"+state.MyPort {
			mapMutex.Unlock()
			log.Println("sent get_kvs_hash to " + addr)
			str, err := json.Marshal(AstatDS.Request{
				Type: AstatDS.GET_KV_HASH,
				IP:   state.MyIP + ":" + state.MyPort,
			})
			if err != nil {
				log.Println(err)
			}
			if conn == nil {
				log.Println("conn is nil")
				UpdateNodeStatus(addr, DEPRECATED)
				continue
			}

			_, err = fmt.Fprintf(conn, string(str)+"\n")
			if err != nil {
				log.Println(err)
				UpdateNodeStatus(addr, DEPRECATED)
				continue
			}
			response, _ := bufio.NewReader(conn).ReadString('\n')

			str, _ = json.Marshal(state.KV)
			log.Println(response + " vs " + MD5(str))
			if response != MD5(str) {
				str, _ := json.Marshal(AstatDS.Request{
					Type: AstatDS.GET_KV,
					IP:   state.MyIP + ":" + state.MyPort,
				})

				log.Println("sent get_kvs to " + addr)
				_, err := fmt.Fprintf(conn, string(str)+"\n")
				if err != nil {
					log.Println(err)
					UpdateNodeStatus(addr, DEPRECATED)
					continue
				}
				response, _ := bufio.NewReader(conn).ReadString('\n')
				log.Println(response)
				kv := new(map[string]Value)
				err = json.Unmarshal([]byte(response), &kv)
				if err != nil {
					log.Println(err)
				}
				kvs = append(kvs, *kv)
			}
		} else {
			mapMutex.Unlock()
		}

	}

	UpdateKV(kvs)
}

func UpdateKV(kvs []map[string]Value) {
	mapMutex.Lock()
	for _, kv := range kvs {
		for k, v := range kv {
			if _, ok := state.KV[k]; !ok {
				state.KV[k] = v
			} else if value, _ := state.KV[k]; value.Value != v.Value {
				t, err := time.Parse(time_format, v.Time)
				curT, err := time.Parse(time_format, value.Time)
				if err != nil {
					log.Println(err)
				}
				if t.After(curT) {
					state.KV[k] = v
				}
			}
		}
	}
	mapMutex.Unlock()
}