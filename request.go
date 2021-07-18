package AstatDS

type Request struct {
	Type  string
	Key   string
	Value string
}

const (
	GET_VALUE = "get value request"
	GET_NODES = "get nodes request" //TODO
	PUT_VALUE = "put value request"

	GET_IPS = "get ips request"
	GET_KV  = "get kvs request"
)
