package server

import (
	"encoding/gob"
	//"fmt"
	"net"
	"net/rpc"
	"strconv"
)

type ServerConfig struct {
	Name   string
	Port   int
	Role   string
	Memory int64
}

func StartMQServer(server string, port int, memory int64) error {
	//fmt.Println("StartMQServer - Memory", memory)
	gob.Register(ServerConfig{})
	mqrpc := NewRPC(&ServerConfig{server, port, "Master", memory})
	rpc.Register(mqrpc)
	l, e := net.Listen("tcp", ":"+strconv.Itoa(port))
	defer l.Close()
	if e != nil {
		return e
	}
	for {
		conn, e := l.Accept()
		if e != nil {
			return e
		}
		go func(c net.Conn) {
			defer c.Close()
			rpc.ServeConn(c)
		}(conn)
	}
	return nil
}
