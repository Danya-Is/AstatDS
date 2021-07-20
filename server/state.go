package main

import (
	"AstatDS"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type State struct {
	KV          map[string]Value
	Ips         map[string]Node
	ClusterName string `json:"clusterName"`
	//TODO MyIP       string `json:"myIP"`
	MyClientPort    string `json:"myClientPort"`
	MyPort          string `json:"myPort"`
	DiscoveryIpPort string `json:"discoveryIpPort"`
	NodeName        string `json:"nodeName"`
	StatePath       string `json:"statePath"`
	Hash            string `json:"hash"`
}

type Node struct {
	time   string
	status string
}

type Value struct {
	time  string
	value string
}

const (
	ACTIVATED  = "activated"
	DEPRECATED = "deprecated"

	time_format = "2006-01-02 15:04:05 MST"
)

func (state *State) DiscoveryNodes() {
	conn, _ := net.Dial("tcp", "0.0.0.0:"+state.DiscoveryIpPort)
	str, _ := json.Marshal(AstatDS.Request{
		Type: AstatDS.GET_IPS,
		IP:   "0.0.0.0:" + state.MyPort,
	})
	_, err := fmt.Fprintf(conn, string(str))
	if err != nil {
		panic(err)
	}
	response, _ := bufio.NewReader(conn).ReadString('\n')
	json.Unmarshal([]byte(response), &state.Ips)
	conn.Close()
}

func (state *State) CheckIps() {
	//обход по нодам
	var ips []map[string]Node
	for _, conn := range connections {
		str, _ := json.Marshal(AstatDS.Request{
			Type: AstatDS.GET_IPS_HASH,
		})
		_, err := fmt.Fprintf(conn, string(str))
		if err != nil {
			panic(err)
		}
		response, _ := bufio.NewReader(conn).ReadString('\n')

		str, _ = json.Marshal(state.Ips)
		if response != MD5(str) {
			str, _ := json.Marshal(AstatDS.Request{
				Type: AstatDS.GET_IPS,
				IP:   "0.0.0.0:" + state.MyPort,
			})
			_, err := fmt.Fprintf(conn, string(str))
			if err != nil {
				panic(err)
			}
			response, _ := bufio.NewReader(conn).ReadString('\n')
			ip := new(map[string]Node)
			json.Unmarshal([]byte(response), &ip)
			ips = append(ips, *ip)
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
			} else if n, _ := state.Ips[addr]; n.status != node.status {
				t, err := time.Parse(time_format, node.time)
				curT, err := time.Parse(time_format, n.time)
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
	for _, conn := range connections {
		str, _ := json.Marshal(AstatDS.Request{
			Type: AstatDS.GET_KV_HASH,
		})
		_, err := fmt.Fprintf(conn, string(str))
		if err != nil {
			panic(err)
		}
		response, _ := bufio.NewReader(conn).ReadString('\n')

		str, _ = json.Marshal(state.Ips)
		if response != MD5(str) {
			str, _ := json.Marshal(AstatDS.Request{
				Type: AstatDS.GET_KV,
				IP:   "0.0.0.0:" + state.MyPort,
			})
			_, err := fmt.Fprintf(conn, string(str))
			if err != nil {
				panic(err)
			}
			response, _ := bufio.NewReader(conn).ReadString('\n')
			kv := new(map[string]Value)
			json.Unmarshal([]byte(response), &kv)
			kvs = append(kvs, *kv)
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
			} else if value, _ := state.KV[k]; value.value != v.value {
				t, err := time.Parse(time_format, v.time)
				curT, err := time.Parse(time_format, value.time)
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
