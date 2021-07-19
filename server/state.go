package main

type State struct {
	KV map[string]interface{}
	Ips map[string]interface{}
	ClusterName string `json:"ClusterName"`
	MyClientPort string `json:"myClientPort"`
	MyPort       string `json:"myPort"`
	DiscoveryIpPort string `json:"discoveryIpPort"`
	NodeName  string `json:"nodeName"`
	StatePath string `json:"statePath"`
	Hash string `json:"hash"`
}

func (state *State) DiscoveryNodes() {
	//обходом из discoveryIpPort находим все ноды кластера
	// заполняем IPS
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
