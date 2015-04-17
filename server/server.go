package server

import (
	"net"
	"net/rpc"
	"strconv"
)

type ServerConfig struct {
	Name string
	Port int
	Role string
}

func StartMQServer(server string, port int, mqexit chan string) error {
	mqrpc := NewRPC(&ServerConfig{server, port, "Master"})
	rpc.Register(mqrpc)
	l, e := net.Listen("tcp", ":"+strconv.Itoa(port))
	if e != nil {
		return e
	}
	for {
		conn, e := l.Accept()
		if e != nil {
			mqexit <- e.Error()
		} else {
			go rpc.ServeConn(conn)
		}
	}
}
