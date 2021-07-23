package go_server

import (
	"AstatDS/server"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
)

var connections map[string]Conn

type Conn struct {
	c net.Conn
}

func (conn *Conn) SentRequest(reqName string) error {
	str, _ := json.Marshal(server.Request{
		Type: reqName,
		IP:   state.MyIP + ":" + state.MyPort,
	})

	if conn == nil {
		UpdateNodeStatus(conn.c.RemoteAddr().String(), DEPRECATED)
		return errors.New("connection is nil")
	}

	_, err := fmt.Fprintf(conn.c, string(str)+"\n")
	if err != nil {
		UpdateNodeStatus(conn.c.RemoteAddr().String(), DEPRECATED)
		return err
	}
	return nil
}

func UpdateConnections() {
	if connections == nil {
		connections = make(map[string]Conn)
	}
	for addr := range state.Ips {
		if addr != state.MyIP+":"+state.MyPort {
			mapMutex.Lock()
			if connections[addr].c == nil {
				newConn, err := net.Dial("tcp", addr)
				connections[addr] = Conn{c: newConn}
				if err != nil {
					log.Println(err)
				}
			}
			mapMutex.Unlock()
		}
	}
}
