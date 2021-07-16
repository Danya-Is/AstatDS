package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var db = make(map[string]interface{})

var (
	state State
)

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

func Init() {
	//читаем с диска

	//если стэйт пуст - ничего не делаем
	//если нет - записываем base64 -> json -> struct в стэйт state := State {...}

	//читаем флаги в стэйт

	state.DiscoveryNodes()
}

func CheckIps() {
	//обход по нодам

	//посылаем запрос сервисам GET_IPS
	//обновляем стэйт
}

func CheckKV() {
	//обход по нодам

	//отправляем запрос GET_KV
	//обновляем стэйт
}

func WriteToDisk() {
	//записываем стэйт в файл
}

func Loop() {
	for {
		CheckIps()
		CheckKV()

		//хэш стейта с предыдущим, если изменен :
		WriteToDisk()
	}
}

func main() {

	Init()

	go Loop()
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

	//разобраться с обработкой запросов по техническому порту
	s.ListenAndServe()
}
