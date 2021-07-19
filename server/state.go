package main

import (
	"AstatDS"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

type State struct {
	KV map[string]interface{}

	Ips map[string]Node

	ClusterName string

	myClientPort string
	myPort       string

	discoveryIpPort string

	nodeName  string
	statePath string

	hash string
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
	conn, _ := net.Dial("tcp", state.discoveryIpPort)
	str, _ := json.Marshal(AstatDS.Request{
		Type: AstatDS.GET_IPS,
		IP:   state.myPort,
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
