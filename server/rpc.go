package server

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	. "github.com/eaciit/mq/client"
	. "github.com/eaciit/mq/helper"
	. "github.com/eaciit/mq/msg"
	"io"
	"log"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	secondsToKill int = 10
)

var (
	serverStartIdle time.Time
	isServerIdle    bool = false
)

type Node struct {
	Config    *ServerConfig
	DataCount int64
	DataSize  int64

	client        *MqClient
	StartTime     time.Time
	offlineStart  time.Time
	isOffline     bool
	AllocatedSize int64
}

type MqRPC struct {
	dataMap map[string]int
	items   map[string]MqMsg
	Config  *ServerConfig
	Host    *ServerConfig

	users []MqUser
	nodes []Node
	exit  bool
}

type MqUser struct {
	UserName string
	Password string
	Role     string
}

func (n *Node) ActiveDuration() time.Duration {
	return time.Since(n.StartTime)
}

func NewRPC(cfg *ServerConfig) *MqRPC {
	m := new(MqRPC)
	m.Config = cfg
	m.items = make(map[string]MqMsg)
	m.nodes = []Node{Node{cfg, 0, 0, nil, time.Now(), time.Now(), false, int64(cfg.Memory)}}
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

func (r *MqRPC) findUser(userName string) int {
	found := false
	for i := 0; i < len(r.users) && !found; i++ {
		if r.users[i].UserName == userName {
			return i
		}
	}
	return -1
}

func (r *MqRPC) GetListUsers(key string, result *MqMsg) error {

	listUser := fmt.Sprintf("UserName \t|Password \n")
	for _, u := range r.users {
		listUser = listUser + fmt.Sprintf("%s \t|%s \n", u.UserName, u.Password)
	}

	(*result).Value = listUser
	return nil
}

func (r *MqRPC) RegisterExistingUser(key string, result *MqMsg) error {
	(*result).Value = ""
	file, err := os.Open("user/user.txt")
	if err != nil {
		fmt.Println("Can't open user file!")
		return nil
	}
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		row := scanner.Text()
		rowSplit := strings.Split(row, "|")
		existingUser := MqUser{}
		existingUser.UserName = rowSplit[0]
		existingUser.Password = rowSplit[1]
		existingUser.Role = rowSplit[2]
		r.users = append(r.users, existingUser)
		infoMsg := fmt.Sprintf("Register User: %s", rowSplit[0])
		fmt.Println(infoMsg)
	}
	return nil
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func SaveUserToFile(userName string, password string, role string) error {
	file, err := os.OpenFile("user/user.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open user file")
	}

	n, err := io.WriteString(file, userName+"|"+password+"|"+role+"\n")
	if err != nil {
		errorMsg := fmt.Sprintf("Error saving user to file, %s:%s", n, err)
		Logging(errorMsg, "ERROR")
	}
	file.Close()
	return nil
}

func UpdateUserFile(r *MqRPC) {
	//r := *MqRPC
	file, err := os.OpenFile("user/user.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalln("Failed to open user file")
	}
	fileContent := ""
	for _, u := range r.users {
		fileContent = fileContent + u.UserName + "|" + u.Password + "|" + u.Role + "\n"
	}
	n, err := io.WriteString(file, fileContent)
	if err != nil {
		errorMsg := fmt.Sprintf("Error update user to file, %s:%s", n, err)
		Logging(errorMsg, "ERROR")
	}
	file.Close()
}

func (r *MqRPC) DeleteUser(value MqMsg, result *MqMsg) error {
	UserName := value.Value.(string)
	Users := []MqUser{}
	for _, u := range r.users {
		//listUser = listUser + fmt.Sprintf("%s \t|%s \n", u.UserName, u.Password)
		if u.UserName != UserName {
			Users = append(Users, u)
		}
	}
	r.users = Users
	UpdateUserFile(r)
	(*result).Value = fmt.Sprintf("User:%s has been deleted", UserName)
	return nil
}

func (r *MqRPC) ChangePassword(value MqMsg, result *MqMsg) error {
	UserName := value.Key
	Password := GetMD5Hash(value.Value.(string))
	Role := "admin"
	userFound := false
	for i, u := range r.users {
		//listUser = listUser + fmt.Sprintf("%s \t|%s \n", u.UserName, u.Password)
		if u.UserName == UserName {
			newUser := MqUser{}
			newUser.UserName = UserName
			newUser.Password = Password
			newUser.Role = Role
			r.users[i] = newUser
			userFound = true
		}
	}
	if userFound {
		UpdateUserFile(r)
		result.Value = "Password has changed successfully for user: " + UserName
	} else {
		result.Value = "Cant find user: " + UserName
	}
	return nil
}

func (r *MqRPC) AddUser(value MqMsg, result *MqMsg) error {
	//check existing user
	userName := value.Key
	password := GetMD5Hash(value.Value.(string))
	role := "admin"
	userIndex := r.findUser(userName)
	userFound := userIndex >= 0
	if userFound {
		errorMsg := "Unable to add user:" + userName + ". It is already exist"
		Logging(errorMsg, "ERROR")
		return errors.New(errorMsg)
	}

	newUser := MqUser{}
	newUser.UserName = userName
	newUser.Password = password
	newUser.Role = role
	r.users = append(r.users, newUser)

	//*result = newUser

	//save user to file
	SaveUserToFile(userName, password, role)

	Logging("New User: "+userName+" has been added with password: "+password, "INFO")
	return nil
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
	newNode.AllocatedSize = nodeConfig.Memory /// 1024 / 1024
	newNode.isOffline = false
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

func (r *MqRPC) CheckHealthSlaves(key string, result *MqMsg) error {
	newNodes := []Node{}
	for i, n := range r.nodes {
		//- check health of the slave
		if strings.ToLower(n.Config.Role) == "slave" {
			_, e := NewMqClient(fmt.Sprintf("%s:%d", n.Config.Name, n.Config.Port), 1*time.Second)
			isActive := true
			if e != nil {

				if !n.isOffline {
					//--- set offline to true and start the offline
					n.isOffline = true
					n.offlineStart = time.Now()
					msg := fmt.Sprintf("CHECK HEALTH OF %s:%d, Slave did not response since %s!", n.Config.Name, n.Config.Port, n.offlineStart)
					Logging(msg, "ERROR")
				}

				//errorMsg := fmt.Sprintf("CHECK HEALTH OF %s:%d, Slave did not response since %s!", n.Config.Name, n.Config.Port, n.offlineStart)
				//Logging(errorMsg, "ERROR")

				//-- check timeout to kill
				duration := time.Since(n.offlineStart)
				kill := int(math.Floor(math.Mod(math.Mod(duration.Seconds(), 3600), 60)))
				if kill >= secondsToKill {
					isActive = false
					errorMsg := fmt.Sprintf("SHUTTING DOWN SLAVE %s:%d, after idle more than %d second(s)", n.Config.Name, n.Config.Port, secondsToKill)
					Logging(errorMsg, "INFO")

				}

				//then remove from r.nodes

			} else {
				if n.isOffline {
					errorMsg := fmt.Sprintf("CHECK HEALTH OF %s:%d, Slave is Up Again!", n.Config.Name, n.Config.Port)
					//fmt.Println(errorMsg)
					Logging(errorMsg, "INFO")
				}
				n.isOffline = false
				//errorMsg := fmt.Sprintf("CHECK HEALTH OF %s:%d, FINE!", n.Config.Name, n.Config.Port)
				//fmt.Println(errorMsg)
			}
			if isActive {
				newNodes = append(newNodes, n)
			}
			r.nodes[i] = n
		} else {
			//if master
			newNodes = append(newNodes, n)
		}

	}
	r.nodes = newNodes
	(*result).Value = ""
	return nil
}

func (r *MqRPC) CheckHealthMaster(key string, result *MqMsg) error {
	callbackCmd := ""
	//fmt.Println("cek master")
	_, e := NewMqClient(fmt.Sprintf(key), 1*time.Second)
	if e != nil {
		//fmt.Println(e)
		if !isServerIdle {
			isServerIdle = true
			serverStartIdle = time.Now()
			errorMsg := fmt.Sprintf("CHECK HEALTH MASTER, Master did not response since %s!", serverStartIdle)
			Logging(errorMsg, "ERROR")
		}

		//-- check timeout to kill
		duration := time.Since(serverStartIdle)
		kill := int(math.Floor(math.Mod(math.Mod(duration.Seconds(), 3600), 60)))
		if kill >= secondsToKill {
			errorMsg := fmt.Sprintf("SHUTTING DOWN, after master idle more than %d second(s)", secondsToKill)
			Logging(errorMsg, "INFO")
			callbackCmd = "KILL"
		}

	} else {
		if isServerIdle {
			errorMsg := fmt.Sprintf("CHECK HEALTH OF MASTER, Master is Up Again!")
			//fmt.Println(errorMsg)
			Logging(errorMsg, "INFO")
		}
		isServerIdle = false
		//errorMsg := fmt.Sprintf("CHECK HEALTH OF MASTER, FINE!")
		//fmt.Println(errorMsg)

		callbackCmd = ""
	}
	(*result).Value = callbackCmd
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

			//fmt.Println("Current node Data Size : ", r.nodes[idx].DataSize)
			//fmt.Println("Incoming Data Size : ", int64(buf.Len()))
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
