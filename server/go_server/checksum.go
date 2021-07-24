package main

import (
	"crypto/md5"
	"fmt"
)

func MD5(data []byte) string {
	h := md5.Sum(data)
	return fmt.Sprintf("%x", h)
}
