package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
)

var db = make(map[string]interface{})

func HomeGetHandler(c *gin.Context) {
	// key := c.Params.ByName("key")
	key := c.Query("key")
	value, ok := db[key]
	if ok {
		c.JSON(200, gin.H{
			"key":   key,
			"value": value})
	} else {
		c.JSON(200, gin.H{"key": key, "value": "no value"})
	}
}

func HomePostHandler(c *gin.Context) {
	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	fmt.Println(string(value))
	if err != nil {
		log.Fatal(err)
	}
	var m interface{}
	err = json.Unmarshal(value, &m)
	data := m.(map[string]interface{})
	for k, v := range data {
		for k1 := range db {
			if k1 == k {
				// later should think about it more
				fmt.Println("this key already exists")
				break
			}
		}
		db[k] = v
	}
	c.JSON(200, data)
}

func main() {
	db["k"] = "v"
	router := gin.Default()
	router.GET("/", HomeGetHandler)
	router.PUT("/", HomePostHandler)
	// r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	s := &http.Server{
		Addr:           ":8080", // for communication with client
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
