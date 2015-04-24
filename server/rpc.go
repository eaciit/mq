package server

import (
	"errors"
	"fmt"
	. "github.com/eaciit/mq/client"
	. "github.com/eaciit/mq/helper"
	. "github.com/eaciit/mq/msg"
	"reflect"
	"strconv"
	"time"
)

type Node struct {
	Config    *ServerConfig
	DataCount int64
	DataSize  int64

	client        *MqClient
	StartTime     time.Time
	AllocatedSize int64
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
	return time.Since(n.StartTime)
}

func NewRPC(cfg *ServerConfig) *MqRPC {
	m := new(MqRPC)
	m.Config = cfg
	m.items = make(map[string]MqMsg)
	m.nodes = []Node{Node{cfg, 0, 0, nil, time.Now(), int64(cfg.Memory)}}
	m.Host = cfg
	return m
}

func (r *MqRPC) Ping(key string, result *MqMsg) error {
	//fmt.Println("Allocated memory", r.nodes[0].AllocatedSize)
	pingInfo := fmt.Sprintf("Server is running on port %s\n", strconv.Itoa(r.Config.Port))
	pingInfo = pingInfo + fmt.Sprintf("Node \t| Address \t| Role \t Active \t\t\t| DataCount \t\t\t| DataSize (MB) \t\t\t|  MaxMemorySize (MB)\t\t\t \n")
	for i, n := range r.nodes {
		pingInfo = pingInfo + fmt.Sprintf("Node %d \t| %s:%d \t| %s \t %v \t\t\t| %d \t\t\t| %d \t\t\t | %d \t\t\t \n", i, n.Config.Name, n.Config.Port,
			n.Config.Role,
			n.ActiveDuration(), n.DataCount, (n.DataSize), (n.AllocatedSize/1024/1024))
	}
	(*result).Value = pingInfo
	return nil
}

func (r *MqRPC) Items(key string, result *MqMsg) error {
	buf, e := Encode(r.items)
	result.Value = buf.Bytes()
	return e
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
		errorMsg := "Unable to add node. It is already exist"
		Logging(errorMsg, "ERROR")
		return errors.New(errorMsg)
	}

	//- check the server
	client, e := NewMqClient(fmt.Sprintf("%s:%d", nodeConfig.Name, nodeConfig.Port), 10*time.Second)
	if e != nil {
		errorMsg := fmt.Sprintf("Unable to add node. Could not connect to %s:%d\n", nodeConfig.Name, nodeConfig.Port)
		Logging(errorMsg, "ERROR")
		return errors.New(errorMsg)
	}
	_, e = client.Call("SetSlave", nodeConfig)
	if e != nil {
		errorMsg := "Unable to add node. Could not set node as slave - message: " + e.Error()
		Logging(errorMsg, "ERROR")
		return errors.New(errorMsg)
	}

	newNode := Node{}
	nodeConfig.Role = "Slave"
	newNode.Config = nodeConfig
	newNode.DataCount = 0
	newNode.DataSize = 0
	newNode.client = client
	newNode.StartTime = time.Now()
	newNode.AllocatedSize = nodeConfig.Memory / 1024 / 1024
	r.nodes = append(r.nodes, newNode)
	Logging("New Node has been added successfully", "INFO")
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

func (r *MqRPC) SetLog(value MqMsg, result *MqMsg) error {
	msg := MqMsg{}
	msg.Key = value.Key
	msg.Value = value.Value
	Logging(msg.Value.(string), msg.Key)
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

func (r *MqRPC) GetLogData(value MqMsg, result *MqMsg) error {
	date := value.Key
	time := value.Value.(string)
	logData, _ := GetLogFileData(date, time)
	(*result).Value = logData
	return nil
}

func (r *MqRPC) Set(value MqMsg, result *MqMsg) error {

	// get value msg
	msg := MqMsg{}
	_, e := r.items[value.Key]
	if e == true {
		msg = r.items[value.Key]
	} else {
		msg.Key = value.Key
	}
	msg.Value = value.Value
	buf, _ := Encode(msg.Value)

	// get nodes where ===> r.nodes[j].DataSize+int64(buf.Len()) < r.nodes[j].AllocatedSize
	idxmasuk := make(map[int]int)
	counteridx := 1
	for j := 0; j < len(r.nodes); j++ {
		if r.nodes[j].DataSize+int64(buf.Len()) < r.nodes[j].AllocatedSize {
			// masuk kriteria
			idxmasuk[j] = j
			counteridx++
		}
	}

	// ada node yang available
	if len(idxmasuk) > 0 {
		// get min node berdasarkan idxmasuk (contains)
		var countNd int64
		var idx int

		// pick min Node
		for i := 0; i < len(r.nodes); i++ {
			if _, ok := idxmasuk[i]; ok { // node ada di list map
				if i == 0 {
					//nd = r.nodes[0]
					countNd = r.nodes[0].DataCount
					idx = 0
				} else {
					if countNd > r.nodes[i].DataCount {
						//nd = r.nodes[i]
						countNd = r.nodes[i].DataCount
						idx = i
					}
				}

			} else {
				// all nodes tidak dapat di isikan data, karena maxsize
			}
		}

		g := r.nodes[idx].DataCount
		maxallocate := r.nodes[idx].AllocatedSize

		if maxallocate > (r.nodes[idx].DataSize + int64(buf.Len())) {
			reflect.ValueOf(&r.nodes[idx]).Elem().FieldByName("DataCount").SetInt(g + 1)
			reflect.ValueOf(&r.nodes[idx]).Elem().FieldByName("DataSize").SetInt((r.nodes[idx].DataSize + int64(buf.Len())) / 1024 / 1024)

			fmt.Println("Current node Data Size : ", r.nodes[idx].DataSize)
			fmt.Println("Incoming Data Size : ", int64(buf.Len()))
			fmt.Println("Data have been set to node, ", "Address : ", r.nodes[idx].Config.Name, " Port : ", r.nodes[idx].Config.Port, " Size : ", r.nodes[idx].DataSize, " DataCount : ", r.nodes[idx].DataCount)
			msg.LastAccess = time.Now()
			r.items[value.Key] = msg

			*result = msg

			Logging("New Key : '"+msg.Key+"' has already set with value: '"+msg.Value.(string)+"'", "INFO")
		} else {
			Logging("New Key : '"+msg.Key+"' with value: '"+msg.Value.(string)+"', data cannot be transmit, because of memory Allocation all node reach max limit", "INFO")
		}
	} else {
		Logging("Data cannot be transmit, because of All node reach max limit", "INFO")
	}

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
	Logging("Key : '"+key+"' has been deleted", "INFO")
	return nil
}
