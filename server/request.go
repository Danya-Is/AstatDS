package server

type Request struct {
	Type  string
	IP    string
	Key   string
	Value string
}

const (
	/*клиентские запросы*/
	GET_VALUE = "get value request"
	GET_NODES = "get nodes request"
	PUT_VALUE = "put value request"

	/*запросы между серверами*/
	GET_IPS      = "get ips request"
	GET_IPS_HASH = "get ips hash request"
	GET_KV       = "get kvs request"
	GET_KV_HASH  = "get kvs hash request"
)
