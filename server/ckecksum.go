package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func MD5(state *State) string {
	data, _ := json.Marshal(state)
	h := md5.Sum(data)
	return fmt.Sprintf("%x", h)
}

func FileMD5(path string) string {
	h := md5.New()
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = io.Copy(h, f)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
