package server

import "github.com/mitchellh/mapstructure"

type Node struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

func ConvertToNode(intr interface{}) Node {
	node := Node{}
	err := mapstructure.Decode(intr, &node)
	if err != nil {
		return Node{}
	}
	return node
}
