package server

type State struct {
	KV map[string]interface{}

	Ips map[string]interface{}

	ClusterName string

	myClientPort string
	myPort       string

	discoveryIpPort string

	nodeName  string
	statePath string

	hash string
}

func (state *State) DiscoveryNodes() {
	//обходом из discoveryIpPort находим все ноды кластера
	// заполняем IPS
}
