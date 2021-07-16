package AstatDS

type Request struct {
	Type  string
	Key   string
	Value string
}

const (
	GET_VALUE = "get value request"
	PUT_VALUE = "put value request"

	GET_IPS = "get ips request"
	GET_KV  = "get kvs request"
)
