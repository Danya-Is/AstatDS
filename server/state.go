package main

import (
	"AstatDS"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

type State struct {
	KV              map[string]interface{}
	Ips             map[string]interface{}
	ClusterName     string `json:"ClusterName"`
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

const (
	ACTIVATED  = "activated"
	DEPRECATED = "deprecated"
)

func (state *State) DiscoveryNodes() {
	conn, _ := net.Dial("tcp", state.DiscoveryIpPort)
	str, _ := json.Marshal(AstatDS.Request{
		Type: AstatDS.GET_IPS,
		IP:   state.MyPort,
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

	//посылаем запрос сервисам GET_IPS
	//обновляем стэйт
}

func (state *State) CheckKV() {
	//обход по нодам

	//отправляем запрос GET_KV
	//обновляем стэйт
}
