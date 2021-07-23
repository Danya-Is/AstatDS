package go_server

import (
	"AstatDS"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

type Conn struct {
	c net.Conn
}

func (conn *Conn) SentRequest(reqName string) error {
	str, _ := json.Marshal(AstatDS.Request{
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
