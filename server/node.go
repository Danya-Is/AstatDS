package server

import "github.com/mitchellh/mapstructure"

type Node struct {
	Time   string `json:"time"`
	Status string `json:"status"`
}

func ConvertToNode(intr interface{}) Node {
	node := Node{}
	err := mapstructure.Decode(intr, &node)
	if err != nil {
		return Node{}
	}
	return node
}
