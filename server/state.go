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
