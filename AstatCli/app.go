package main

import (
	"AstatDS/AstatCli/commands"
	"io"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	app := cli.NewApp()
	app.Name = "Astat"
	app.Usage = "write"
	//app.Action = func(c *cli.Context) error {
	//	fmt.Println("Hello friend!")
	//	serv := "localhost:8080"
	//	conn, _ := net.Dial("tcp", serv) // открываем TCP-соединение к серверу
	//	go copyTo(os.Stdout, conn)       // читаем из сокета в stdout
	//	copyTo(conn, os.Stdin)           // пишем в сокет из stdin
	//	return nil
	//}
	app.Commands = []*cli.Command{
		commands.NewSetConfigCommand(),
	}
	app.Run(os.Args)
}

func copyTo(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
