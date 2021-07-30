package main

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

func (conn *Conn) SentRequest(reqName string, addr string) error {
	str, _ := json.Marshal(server.Request{
		Type: reqName,
		IP:   state.MyIP + ":" + state.MyPort,
	})

	if conn.c == nil {
		UpdateNodeStatus(addr, DEPRECATED)
		return errors.New("connection is nil")
	}

	_, err := fmt.Fprintf(conn.c, string(str)+"\n")
	if err != nil {
		UpdateNodeStatus(addr, DEPRECATED)
		return err
	}
	return nil
}

func UpdateConnections() {
	if connections == nil {
		connections = make(map[string]Conn)
	}
	for i := 0; i < len(Ips.Keys()); i++ {
		ip := fmt.Sprint(Ips.Keys()[i])
		if ip != state.MyIP+":"+state.MyPort {
			mapMutex.Lock()
			str, _ := Ips.Get(ip)
			mapMutex.Unlock()
			con := server.ConvertToNode(str)
			if connections[ip].c == nil || con.Status == DEPRECATED {
				newConn, err := net.Dial("tcp", ip)
				connections[ip] = Conn{c: newConn}
				if err != nil {
					log.Println(err)
				} else {
					UpdateNodeStatus(ip, ACTIVATED)
				}
			}
		}
	}
}
