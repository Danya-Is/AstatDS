package main

import (
	"crypto/md5"
	"fmt"
)

func MD5(data []byte) string {
	h := md5.Sum(data)
	return fmt.Sprintf("%x", h)
}

/*func main() {
	fmt.Println(MD5([]byte("{\"C\":{\"time\":\"2021-07-24 01:54:25 MSK\",\"value\":\"C\"},\"I\":{\"time\":\"2021-07-24 01:54:34 MSK\",\"value\":\"I\"},\"aaa\":{\"time\":\"2021-07-24 01:54:11 MSK\",\"value\":\"aaa\"},\"hello2\":{\"time\":\"2021-07-23 15:36:58 MSK\",\"value\":\"hello2 test server\"},\"hello3\":{\"time\":\"2021-07-23 15:36:58 MSK\",\"value\":\"hello3 test server\"},\"hello4\":{\"time\":\"2021-07-23 15:36:58 MSK\",\"value\":\"hello4 test server\"},\"key\":{\"time\":\"2021-07-24 01:54:07 MSK\",\"value\":\"value\"},\"zzz\":{\"time\":\"2021-07-24 01:54:09 MSK\",\"value\":\"zzz\"}}")))
}*/
