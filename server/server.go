package server

import (
	"encoding/gob"
	"net"
	"net/rpc"
	"strconv"
)

type ServerConfig struct {
	Name string
	Port int
	Role string
}

func StartMQServer(server string, port int) error {
	gob.Register(ServerConfig{})
	mqrpc := NewRPC(&ServerConfig{server, port, "Master"})
	rpc.Register(mqrpc)
	l, e := net.Listen("tcp", ":"+strconv.Itoa(port))
	if e != nil {
		return e
	}
	for {
		conn, e := l.Accept()
		if e != nil {
			return e
		}
		if conn != nil {
			go func(c net.Conn) {
				defer c.Close()
				rpc.ServeConn(c)
			}(conn)
		}
	}
}
