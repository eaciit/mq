package server

import (
	"errors"
	"fmt"
	. "github.com/eaciit/mq/msg"
	. "strconv"
	"time"
)

type Node struct {
	Config    *ServerConfig
	DataCount int64
	DataSize  int64
	NodeRole  string
}

type MqRPC struct {
	items  map[string]MqMsg
	Config *ServerConfig
	Host   *ServerConfig

	nodeRole  string
	nodes     []Node
	exit      bool
	startTime time.Time
}

func NewRPC(cfg *ServerConfig) *MqRPC {
	m := new(MqRPC)
	m.Config = cfg
	m.items = make(map[string]MqMsg)
	m.nodes = make([]Node, 0)
	m.startTime = time.Now()

	m.Host = cfg
	m.nodeRole = "Master"
	return m
}

func (r *MqRPC) Ping(key string, result *MqMsg) error {
	runDuration := time.Since(r.startTime)
	(*result).Value = fmt.Sprintf("Server is running on port %s  since %v (%v) \n", Itoa(r.Config.Port), r.startTime, runDuration)
	return nil
}

func (r *MqRPC) findNode(serverName string, port int) (int, Node) {
	found := false
	for i := 0; i < len(r.nodes) || !found; i++ {
		if r.nodes[i].Config.Name == serverName && r.nodes[i].Config.Port == port {
			return i, r.nodes[i]
		}
	}
	return -1, Node{}
}

func (r *MqRPC) AddNode(nodeConfig *ServerConfig, result *MqMsg) error {
	nodeIndex, _ := r.findNode(nodeConfig.Name, nodeConfig.Port)
	nodeFound := nodeIndex >= 0
	if nodeFound {
		return errors.New("Unable to add slave. It is already exist")
	}
	newNode := Node{}
	newNode.Config = nodeConfig
	newNode.DataCount = 0
	newNode.DataSize = 0
	r.nodes = append(r.nodes, newNode)
	return nil
}

func (r *MqRPC) Kill(key string, result *MqMsg) error {
	r.exit = true
	(*result).Value = ""
	return nil
}

func (r *MqRPC) GetLog(key time.Time, result *MqMsg) error {
	if r.exit {
		(*result).Value = fmt.Sprintf("Received EXIT command at %v \n", time.Now())
	} else {
		(*result).Value = ""
	}
	return nil
}
