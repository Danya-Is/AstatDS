package main

import (
	//"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"io"
	"log"
	"net"
	"os"
)

func main() {

	app := cli.NewApp()
	app.Name = "Astat"
	app.Usage = "write"
	app.Action = func(c *cli.Context) error {
		fmt.Println("Hello friend!")
		serv := "localhost:8080"
		conn, _ := net.Dial("tcp", serv) // открываем TCP-соединение к серверу
		go copyTo(os.Stdout, conn)       // читаем из сокета в stdout
		copyTo(conn, os.Stdin)           // пишем в сокет из stdin
		return nil
	}
	// app.Commands = []*cli.Command{
	// 	{
	// 		Name:  "write",
	// 		Usage: "write a json object to server",
	// 		Action: func(c *cli.Context) {
	// 			//jsonObj, err := json.Marhsal({c.Args.First() : c.Args.Second()})
	// 			serv := "localhost:8080"
	// 			conn, _ := net.Dial("tcp", serv) // открываем TCP-соединение к серверу
	// 			go copyTo(os.Stdout, conn)       // читаем из сокета в stdout
	// 			copyTo(conn, os.Stdin)           // пишем в сокет из stdin

	// 		}}}
	app.Run(os.Args)
}

func copyTo(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
