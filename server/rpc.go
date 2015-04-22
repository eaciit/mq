package server

import (
	"errors"
	"fmt"
	. "github.com/eaciit/mq/client"
	. "github.com/eaciit/mq/helper"
	. "github.com/eaciit/mq/msg"
	"strconv"
	"time"
)

type Node struct {
	Config    *ServerConfig
	DataCount int64
	DataSize  int64

	client    *MqClient
	startTime time.Time
}

type MqRPC struct {
	dataMap map[string]int
	items   map[string]MqMsg
	Config  *ServerConfig
	Host    *ServerConfig

	nodes []Node
	exit  bool
}

func (n *Node) ActiveDuration() time.Duration {
	return time.Since(n.startTime)
}

func NewRPC(cfg *ServerConfig) *MqRPC {
	m := new(MqRPC)
	m.Config = cfg
	m.items = make(map[string]MqMsg)
	m.nodes = []Node{Node{cfg, 0, 0, nil, time.Now()}}
	m.Host = cfg
	return m
}

func (r *MqRPC) Ping(key string, result *MqMsg) error {
	pingInfo := fmt.Sprintf("Server is running on port %s\n", strconv.Itoa(r.Config.Port))
	pingInfo = pingInfo + fmt.Sprintf("Node \t| Address \t| Role \t Active \t\t\t| Data# \t\t\t| Data(MB) \n")
	for i, n := range r.nodes {
		pingInfo = pingInfo + fmt.Sprintf("Node %d \t| %s:%d \t| %s \t %v \t\t\t| %d \t\t\t| %d \n", i, n.Config.Name, n.Config.Port,
			n.Config.Role,
			n.ActiveDuration(), n.DataCount, (n.DataSize/1024/1024))
	}
	(*result).Value = pingInfo
	return nil
}

func (r *MqRPC) Nodes(key string, result *MqMsg) error {
	buf, e := Encode(r.nodes)
	result.Value = buf.Bytes()
	return e
}

func (r *MqRPC) findNode(serverName string, port int) (int, Node) {
	found := false
	for i := 0; i < len(r.nodes) && !found; i++ {
		if r.nodes[i].Config.Name == serverName && r.nodes[i].Config.Port == port {
			return i, r.nodes[i]
		}
	}
	return -1, Node{}
}

func (r *MqRPC) AddNode(nodeConfig *ServerConfig, result *MqMsg) error {
	//-- is server exist
	nodeIndex, _ := r.findNode(nodeConfig.Name, nodeConfig.Port)
	nodeFound := nodeIndex >= 0
	if nodeFound {
		return errors.New("Unable to add node. It is already exist")
	}

	//- check the server
	client, e := NewMqClient(fmt.Sprintf("%s:%d", nodeConfig.Name, nodeConfig.Port), 10*time.Second)
	if e != nil {
		return errors.New(fmt.Sprintf("Unable to add node. Could not connect to %s:%d\n", nodeConfig.Name, nodeConfig.Port))
	}
	_, e = client.Call("SetSlave", nodeConfig)
	if e != nil {
		return errors.New("Unable to add node. Could not set node as slave - message: " + e.Error())
	}

	newNode := Node{}
	nodeConfig.Role = "Slave"
	newNode.Config = nodeConfig
	newNode.DataCount = 0
	newNode.DataSize = 0
	newNode.client = client
	newNode.startTime = time.Now()
	r.nodes = append(r.nodes, newNode)
	return nil
}

func (r *MqRPC) GetConfig(key string, result *MqMsg) error {
	result.Value = *r.Config
	return nil
}

func (r *MqRPC) SetSlave(config *ServerConfig, result *MqMsg) error {
	r.Config.Role = "Slave"
	r.Host = config
	r.nodes = []Node{}
	return nil
}

func (r *MqRPC) Kill(key string, result *MqMsg) error {
	for _, n := range r.nodes {
		if n.Config.Role != "Master" {
			n.client.Call("Kill", "")
		}
	}
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

func (r *MqRPC) Set(value MqMsg, result *MqMsg) error {
	msg := MqMsg{}
	_, e := r.items[value.Key]
	if e == true {
		msg = r.items[value.Key]
	} else {
		msg.Key = value.Key
	}
	msg.Value = value.Value
	msg.LastAccess = time.Now()
	r.items[value.Key] = msg

	*result = msg
	return nil
}

func (r *MqRPC) Get(key string, result *MqMsg) error {
	v, e := r.items[key]
	if e == false {
		return errors.New("Data for key " + key + " is not exist")
	}
	*result = v
	return nil
}

func (r *MqRPC) Delete(key string, result *MqMsg) error {
	_, e := r.items[key]
	if e == true {
		delete(r.items, key)
	}
	return nil
}
