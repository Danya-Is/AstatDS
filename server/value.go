package server

import "github.com/mitchellh/mapstructure"

type Value struct {
	Time  string `json:"time"`
	Value string `json:"value"`
}

func ConvertToValue(intr interface{}) Value {
	v := Value{}
	err := mapstructure.Decode(intr, &v)
	if err != nil {
		return Value{}
	}
	return v
}
