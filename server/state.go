package main

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
	fmt.Println("str writed")
	if err != nil {
		log.Fatal(err)
	}
	response, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println(response)
	json.Unmarshal([]byte(response), &state.Ips)
	fmt.Println("nodes discovered")
	fmt.Println(state)
	conn.Close()
}

func (state *State) CheckIps() {
	//обход по нодам
	var ips []map[string]Node
	for addr, conn := range connections {
		if state.Ips[addr].Status == ACTIVATED {
			str, _ := json.Marshal(AstatDS.Request{
				Type: AstatDS.GET_IPS_HASH,
			})
			if conn == nil {
				state.Ips[addr] = Node{
					Status: DEPRECATED,
					Time:   time.Now().Format(time_format),
				}
				continue
			}
			_, err := fmt.Fprintf(conn, string(str)+"\n")
			if err != nil {
				log.Println(err)
				state.Ips[addr] = Node{
					Status: DEPRECATED,
					Time:   time.Now().Format(time_format),
				}
				continue
			}
			response, _ := bufio.NewReader(conn).ReadString('\n')

			str, _ = json.Marshal(state.Ips)
			if response != MD5(str) {
				str, _ := json.Marshal(AstatDS.Request{
					Type: AstatDS.GET_IPS,
					IP:   state.MyIP + ":" + state.MyPort,
				})
				_, err := fmt.Fprintf(conn, string(str)+"\n")
				if err != nil {
					log.Println(err)
					state.Ips[addr] = Node{
						Status: DEPRECATED,
						Time:   time.Now().Format(time_format),
					}
					continue
				}
				response, _ := bufio.NewReader(conn).ReadString('\n')
				ip := new(map[string]Node)
				json.Unmarshal([]byte(response), &ip)
				ips = append(ips, *ip)
			}
		}
	}
	//посылаем запрос сервисам GET_IPS

	UpdateIps(ips)
	//обновляем стэйт
}

func UpdateIps(ips []map[string]Node) {
	for _, m := range ips {
		for addr, node := range m {
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
		}
	}
}

func (state *State) CheckKV() {
	//обход по нодам
	var kvs []map[string]Value
	for addr, conn := range connections {
		if state.Ips[addr].Status == ACTIVATED && addr != state.MyIP+":"+state.MyPort {
			//conn, err := net.Dial("tcp", addr)
			//if err != nil {
			//	log.Println(err)
			//	continue
			//}
			str, err := json.Marshal(AstatDS.Request{
				Type: AstatDS.GET_KV_HASH,
			})
			if err != nil {
				log.Println(err)
			}
			if conn == nil {
				state.Ips[addr] = Node{
					Status: DEPRECATED,
					Time:   time.Now().Format(time_format),
				}
				continue
			}
			_, err = fmt.Fprintf(conn, string(str)+"\n")
			if err != nil {
				log.Println(err)
				state.Ips[addr] = Node{
					Status: DEPRECATED,
					Time:   time.Now().Format(time_format),
				}
				continue
			}
			response, _ := bufio.NewReader(conn).ReadString('\n')

			str, _ = json.Marshal(state.KV)
			//fmt.Println(response + " vs " + MD5(str))
			if response != MD5(str) {
				str, _ := json.Marshal(AstatDS.Request{
					Type: AstatDS.GET_KV,
					IP:   state.MyIP + ":" + state.MyPort,
				})
				fmt.Println(str)
				_, err := fmt.Fprintf(conn, string(str)+"\n")
				if err != nil {
					log.Println(err)
					state.Ips[addr] = Node{
						Status: DEPRECATED,
						Time:   time.Now().Format(time_format),
					}
					continue
				}
				response, _ := bufio.NewReader(conn).ReadString('\n')
				kv := new(map[string]Value)
				json.Unmarshal([]byte(response), &kv)
				kvs = append(kvs, *kv)
			}
		}
	}
	//отправляем запрос GET_KV

	UpdateKV(kvs)
	//обновляем стэйт
}

func UpdateKV(kvs []map[string]Value) {
	for _, kv := range kvs {
		for k, v := range kv {
			if _, ok := state.KV[k]; !ok {
				state.KV[k] = v
			} else if value, _ := state.KV[k]; value.Value != v.Value {
				t, err := time.Parse(time_format, v.Time)
				curT, err := time.Parse(time_format, value.Time)
				if err != nil {
					panic(err)
				}
				if t.After(curT) {
					state.KV[k] = v
				}
			}
		}
	}
}
