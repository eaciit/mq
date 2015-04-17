package server

import (
	"fmt"
	. "github.com/eaciit/mq/msg"
	. "strconv"
)

type MqRPC struct {
	items map[string]MqMsg

	Config *ServerConfig
}

func NewRPC(cfg *ServerConfig) *MqRPC {
	m := new(MqRPC)
	m.Config = cfg
	m.items = make(map[string]MqMsg)
	return m
}

func (r *MqRPC) Info(key string, result *MqMsg) error {
	(*result).Value = fmt.Sprintf("Server is running on port %s", Itoa(r.Config.Port))
	return nil
}

func (r *MqRPC) Get(key string, result *string) error {
	return nil
}
