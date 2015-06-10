package server

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/eaciit/mq/client"
	. "github.com/eaciit/mq/helper"
	. "github.com/eaciit/mq/msg"
)

const (
	secondsToKill int = 10
)

var (
	serverStartIdle time.Time
	isServerIdle    bool = false
)

type Pair struct {
	First, Second interface{}
}

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
	dataMap        map[string]int
	items          map[string]MqMsg
	tables         map[string]MqTable
	Config         *ServerConfig
	Host           *ServerConfig
	deadNodesCount int

	users   []MqUser
	nodes   []Node
	mirrors []Node
	exit    bool
}

type Table struct {
	Key   string
	Value string
	Owner string
}

type MqUser struct {
	UserName    string
	Password    string
	Role        string
	DateCreated time.Time
}

func (n *Node) ActiveDuration() time.Duration {
	return time.Since(n.StartTime)
}

func NewRPC(cfg *ServerConfig) *MqRPC {
	m := new(MqRPC)
	m.dataMap = make(map[string]int)
	m.Config = cfg
	m.items = make(map[string]MqMsg)
	m.tables = make(map[string]MqTable)
	m.nodes = []Node{Node{cfg, 0, 0, nil, time.Now(), time.Now(), false, int64(cfg.Memory)}}
	m.mirrors = []Node{}
	m.Host = cfg
	return m
}

func (r *MqRPC) Ping(key string, result *MqMsg) error {
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

func (r *MqRPC) Users(key string, result *MqMsg) error {
	buf, e := Encode(r.users)
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

	// Adding searching for existing mirror
	for i := 0; i < len(r.mirrors) && !found; i++ {
		if r.mirrors[i].Config.Name == serverName && r.mirrors[i].Config.Port == port {
			return i, r.mirrors[i]
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
		layout := "Mon, 01/02/06, 03:04PM"
		t, _ := time.Parse(layout, rowSplit[3])
		existingUser.DateCreated = t
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
		fileContent = fileContent + fmt.Sprintf("%s|%s|%s|%s\n", u.UserName, u.Password, u.Role, u.DateCreated)
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
			newUser.DateCreated = r.users[i].DateCreated
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

func (r *MqRPC) ClientLogin(value MqMsg, result *MqMsg) error {
	UserName := value.Key
	Password := GetMD5Hash(value.Value.(string))
	Role := ""
	userFound := false

	if UserName == "root" && Password == GetMD5Hash("Password.1") {
		userFound = true
		Role = "root"
	} else {
		for _, u := range r.users {
			//listUser = listUser + fmt.Sprintf("%s \t|%s \n", u.UserName, u.Password)
			if u.UserName == UserName {
				if u.Password == Password {
					userFound = true
					Role = u.Role
				}
			}
		}
	}
	if userFound {
		result.Value = Role
	} else {
		result.Value = "0"
	}
	return nil
}

func (r *MqRPC) AddUser(value MqMsg, result *MqMsg) error {
	//check existing user
	splitKey := strings.Split(value.Key, "|")
	userName := splitKey[0]
	role := splitKey[1]
	if role == "" {
		role = "admin"
	}
	password := GetMD5Hash(value.Value.(string))
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
	newUser.DateCreated = time.Now()
	r.users = append(r.users, newUser)

	// *result = newUser

	//save user to file
	UpdateUserFile(r)

	Logging("New User: "+userName+" has been added with password: "+password, "INFO")
	return nil
}

func (r *MqRPC) AddNode(nodeConfig *ServerConfig, result *MqMsg) error {
	//-- is server exist
	nodeIndex, _ := r.findNode(nodeConfig.Name, nodeConfig.Port)
	nodeFound := nodeIndex >= 0
	if nodeFound {
		errorMsg := fmt.Sprintf("Unable to add node %s:%s. It is already exist", nodeConfig.Name, nodeConfig.Port)
		Logging(errorMsg, "ERROR")
		return errors.New(errorMsg)
	}

	//- check the slave server
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

	if r.deadNodesCount > 0 {
		// There is dead node in master metadata
		// Get meta for dead node
		lostMeta := []string{}
		for index, item := range r.dataMap {
			if item == -r.deadNodesCount {
				lostMeta = append(lostMeta, index)
			}
		}

		lastNodeIndex := len(r.nodes) - 1
		err := r.copyDataFromMirrorToNode(lastNodeIndex, lostMeta)
		if err == nil {
			// Copying data succes, update the dataMap
			for index, item := range r.dataMap {
				if item == -r.deadNodesCount {
					r.dataMap[index] = lastNodeIndex
				}
			}
		}
	}

	return nil
}

// TODO: in AddMirror and AddNode there are some similar code, so merge them!
func (r *MqRPC) AddMirror(mirrorConfig *ServerConfig, result *MqMsg) error {
	nodeIndex, _ := r.findNode(mirrorConfig.Name, mirrorConfig.Port)
	nodeFound := nodeIndex >= 0
	if nodeFound {
		errorMsg := fmt.Sprintf("Unable to add mirror %s:%s. It is already exist", mirrorConfig.Name, mirrorConfig.Port)
		Logging(errorMsg, "ERROR")
		return errors.New(errorMsg)
	}

	//- check the mirror server
	client, e := NewMqClient(fmt.Sprintf("%s:%d", mirrorConfig.Name, mirrorConfig.Port), 10*time.Second)
	if e != nil {
		errorMsg := fmt.Sprintf("Unable to add mirror. Could not connect to %s:%d\n", mirrorConfig.Name, mirrorConfig.Port)
		Logging(errorMsg, "ERROR")
		return errors.New(errorMsg)
	}
	_, e = client.Call("SetMirror", mirrorConfig)
	if e != nil {
		errorMsg := "Unable to add node. Could not set node as mirror - message: " + e.Error()
		Logging(errorMsg, "ERROR")
		return errors.New(errorMsg)
	}

	newNode := Node{}
	mirrorConfig.Role = "Mirror"
	newNode.Config = mirrorConfig
	newNode.DataCount = 0
	newNode.DataSize = 0
	newNode.client = client
	newNode.StartTime = time.Now()
	newNode.AllocatedSize = mirrorConfig.Memory /// 1024 / 1024
	newNode.isOffline = false
	r.mirrors = append(r.mirrors, newNode)
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

func (r *MqRPC) SetMirror(config *ServerConfig, result *MqMsg) error {
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
	// fmt.Println(len(r.items))
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
					// Kill node
					// Updating dead nodes count
					r.deadNodesCount += 1

					// Update dataMap before killing node
					for index, item := range r.dataMap {
						if item == i {
							r.dataMap[index] = -r.deadNodesCount
						}
					}

					isActive = false
					errorMsg := fmt.Sprintf("SHUTTING DOWN SLAVE %s:%d, after idle more than %d second(s)", n.Config.Name, n.Config.Port, secondsToKill)
					Logging(errorMsg, "INFO")
				}

				//then remove from r.nodes

			} else {
				if n.isOffline {
					r.checkReconnectedNode(i)
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

				// Previous node is dead
				if len(newNodes) == i {
					// Updating data map
					for index, item := range r.dataMap {
						if item == i {
							// Update index to Previous node because dead node already deleted
							r.dataMap[index] = i - (i - len(newNodes)) - 1
						}
					}
				}

			}
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
	// fmt.Println(len(r.items))
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

func parseValue(value string, result *MqMsg) error {

	valsplit := strings.Split(value, "|")
	for i := 0; i <= len(valsplit); i++ {
		if i == 0 {

		} else {
			if strings.Contains(strings.ToLower(valsplit[i]), "owner") {
				result.Owner = strings.Split(valsplit[i], "=")[1]
			}
			if strings.Contains(strings.ToLower(valsplit[i]), "duration") {
				//result.Duration = int64(strings.Split(valsplit[i], "=")[1])
			}
			if strings.Contains(strings.ToLower(valsplit[i]), "table") {
				result.Table = strings.Split(valsplit[i], "=")[1]
			}
			if strings.Contains(strings.ToLower(valsplit[i]), "permission") {
				result.Permission = strings.Split(valsplit[i], "=")[1]
			}

		}
	}
	fmt.Println("parseValue", result)
	//result.
	return nil
}

func (r *MqRPC) checkReconnectedNode(nodeIndex int) error {
	n := r.nodes[nodeIndex]
	client, err := NewMqClient(fmt.Sprintf("%s:%d", n.Config.Name, n.Config.Port), 1*time.Second)
	if err != nil {
		errorMsg := fmt.Sprintf("Unable connect to node %s:%d\n", n.Config.Name, n.Config.Port)
		Logging(errorMsg, "ERROR")
		return errors.New(errorMsg)
	}

	args := []string{}
	for key, value := range r.dataMap {
		if value == nodeIndex {
			args = append(args, key)
		}
	}

	lostMeta := []string{}
	err = client.CallDirect("CheckData", args, &lostMeta)
	if err != nil {
		errorMsg := fmt.Sprintf("Unable connect to node %s:%d\n", n.Config.Name, n.Config.Port)
		Logging(errorMsg, "ERROR")
		return errors.New(errorMsg)
	}

	if len(lostMeta) > 0 {
		r.copyDataFromMirrorToNode(nodeIndex, lostMeta)
	} else {
		Logging("All data still exist", "INFO")
	}

	return nil
}

func (r *MqRPC) copyDataFromMirrorToNode(nodeIndex int, lostMeta []string) error {
	Logging("Data lost, copying data from mirror", "INFO")
	if len(r.mirrors) > 0 {

		for _, mirror := range r.mirrors {
			//Getting data from mirror
			mirrorClient, err := NewMqClient(fmt.Sprintf("%s:%d", mirror.Config.Name, mirror.Config.Port), 1*time.Second)
			if err != nil {
				errorMsg := fmt.Sprintf("Unable connect to mirror %s:%d\n", mirror.Config.Name, mirror.Config.Port)
				Logging(errorMsg, "ERROR")
			}

			result := false
			args := Pair{r.nodes[nodeIndex].Config, lostMeta}
			err = mirrorClient.CallDirect("FindAndSendItems", args, &result)
			if err != nil {
				errorMsg := fmt.Sprintf("Unable send command to mirror %s:%d, Error: %s \n", mirror.Config.Name, mirror.Config.Port, err.Error())
				Logging(errorMsg, "ERROR")
			}

			Logging("Data successfully retrivied by node from mirror", "INFO")
		}

	} else {
		errorMsg := "Cant copy data to new node, no existing mirror connected"
		Logging(errorMsg, "ERROR")
		return errors.New(errorMsg)
	}

	return nil
}

// Check existing data with master metadata
func (r *MqRPC) CheckData(args []string, result *[]string) error {
	for _, key := range args {
		_, exist := r.items[key]
		if !exist {
			args = append(*result, key)
		}
	}

	return nil
}

func (r *MqRPC) FindAndSendItems(args Pair, result *bool) error {
	nodeConfig := args.First.(ServerConfig)
	lostMeta := args.Second.([]string)

	selectedItem := make(map[string]MqMsg)
	for _, key := range lostMeta {
		selectedItem[key] = r.items[key]
	}

	client, err := NewMqClient(fmt.Sprintf("%s:%d", nodeConfig.Name, nodeConfig.Port), 1*time.Second)
	if err != nil {
		errorMsg := fmt.Sprintf("Unable connect to node %s:%d\n", nodeConfig.Name, nodeConfig.Port)
		Logging(errorMsg, "ERROR")
		return errors.New(errorMsg)
	}

	err = client.CallDirect("RetrieveDatas", selectedItem, result)
	if err != nil {
		errorMsg := fmt.Sprintf("Unable send data to slaves from mirror %s:%d\n", nodeConfig.Name, nodeConfig.Port)
		Logging(errorMsg, "ERROR")
		return errors.New(errorMsg)
	}

	return nil
}

func (r *MqRPC) RetrieveDatas(datas map[string]MqMsg, result *bool) error {
	for key, data := range datas {
		r.items[key] = data
	}

	*result = true
	return nil
}

func (r *MqRPC) Set(value MqMsg, result *MqMsg) error {
	msg := MqMsg{}
	// check if data already in items
	_, exist := r.items[value.Key]
	if exist {
		msg = r.items[value.Key]
	} else {
		msg.Key = value.Key
	}
	msg.Value = value.Value

	buf, _ := Encode(msg.Value)

	// Search for available node
	idxmasuk := make(map[int]int)
	counteridx := 1
	for j := 0; j < len(r.nodes); j++ {
		if r.nodes[j].DataSize+int64(buf.Len()) < r.nodes[j].AllocatedSize {
			idxmasuk[j] = j
			counteridx++
		}
	}

	if len(idxmasuk) > 0 {
		var countNd int64
		var idx int

		// Pick min node
		for i := 0; i < len(r.nodes); i++ {
			if _, ok := idxmasuk[i]; ok {
				if i == 0 {
					countNd = r.nodes[0].DataCount
					idx = 0
				} else {
					if countNd > r.nodes[i].DataCount {
						countNd = r.nodes[i].DataCount
						idx = i
					}
				}
			}
		}

		// Store data to all existing mirror
		for key, mirror := range r.mirrors {
			r.mirrors[key].DataCount += 1
			r.mirrors[key].DataSize += int64(buf.Len()) / 1024 / 1024

			client, e := NewMqClient(fmt.Sprintf("%s:%d", mirror.Config.Name, mirror.Config.Port), 10*time.Second)
			if e != nil {
				errorMsg := fmt.Sprintf("Unable connect to node %s:%d\n", r.nodes[idx].Config.Name, r.nodes[idx].Config.Port)
				Logging(errorMsg, "ERROR")
				return errors.New(errorMsg)
			}

			client.Call("SetItem", msg)
			if e != nil {
				errorMsg := fmt.Sprintf("Unable to set data to node : %s", e.Error())
				return errors.New(errorMsg)
			}

			fmt.Printf("Data has been mirrored to Address: %s:%d, Size: %d DataCount: %d\n", mirror.Config.Name, mirror.Config.Port, r.mirrors[key].DataSize, r.mirrors[key].DataCount)
		}

		if r.nodes[idx].AllocatedSize > (r.nodes[idx].DataSize + int64(buf.Len())) {
			r.nodes[idx].DataCount += 1
			r.nodes[idx].DataSize += int64(buf.Len()) / 1024 / 1024

			fmt.Println("Data has been set to node, ", "Address : ", r.nodes[idx].Config.Name, " Port : ", r.nodes[idx].Config.Port, " Size : ", r.nodes[idx].DataSize, " DataCount : ", r.nodes[idx].DataCount)
			msg.LastAccess = time.Now()
			msg.SetDefaults(&msg)

			// Decode data
			valsplit := strings.Split(value.Value.(string), "|")
			for i := 0; i < len(valsplit); i++ {
				field := strings.ToLower(strings.Split(valsplit[i], "=")[0])
				if strings.TrimSpace(field) == "owner" {
					msg.Owner = strings.TrimSpace(strings.Split(valsplit[i], "=")[1])
					msg.Owner = strings.Trim(msg.Owner, "\"")
				}
				if strings.TrimSpace(field) == "duration" {
					x, _ := strconv.ParseInt(strings.Split(valsplit[i], "=")[1], 0, 64)
					msg.Duration = x //strings.Split(valsplit[i], "=")[1].(int64))
				}
				if strings.TrimSpace(field) == "table" {
					msg.Table = strings.TrimSpace(strings.Split(valsplit[i], "=")[1])
					msg.Table = strings.Trim(msg.Table, "\"")
				}
				if strings.TrimSpace(field) == "permission" {
					msg.Permission = strings.TrimSpace(strings.Split(valsplit[i], "=")[1])
					msg.Permission = strings.Trim(msg.Permission, "\"")
				}
			}
			msg.Key = value.Key
			*result = msg

			// Set item to selected node
			client, e := NewMqClient(fmt.Sprintf("%s:%d", r.nodes[idx].Config.Name, r.nodes[idx].Config.Port), 10*time.Second)
			if e != nil {
				errorMsg := fmt.Sprintf("Unable connect to node %s:%d\n", r.nodes[idx].Config.Name, r.nodes[idx].Config.Port)
				Logging(errorMsg, "ERROR")
				return errors.New(errorMsg)
			}

			client.Call("SetItem", msg)
			if e != nil {
				errorMsg := fmt.Sprintf("Unable to set data to node : %s", e.Error())
				return errors.New(errorMsg)
			}

			r.dataMap[msg.Key] = idx
			r.setTableProperties(msg)
			Logging("New Key : '"+msg.Key+"' has already set with value: '"+msg.Value.(string)+"'", "INFO")
		} else {
			Logging("New Key : '"+msg.Key+"' with value: '"+msg.Value.(string)+"', data cannot be transmit, because of memory Allocation all node reach max limit", "INFO")
		}
	} else {
		Logging("Data cannot be transmit, because of All node reach max limit", "INFO")
	}

	return nil
}

func (r *MqRPC) SetItem(data MqMsg, result *MqMsg) error {
	r.items[data.Key] = data
	*result = data
	return nil
}

func (r *MqRPC) Inc(key map[string]interface{}, result *MqMsg) error {
	k := key["key"]
	data := key["data"]
	v, e := r.items[k.(string)]
	if e == false {
		return errors.New("Data for key  is not exist")
	} else {
		v.Value = data
		r.items[k.(string)] = v
	}
	return nil
}

func (r *MqRPC) GetItem(key string, result *MqMsg) error {
	v, e := r.items[key]
	if e == false {
		return errors.New("Data for key " + key + " is not exist")
	}
	*result = v
	return nil
}

func (r *MqRPC) Get(key string, result *MqMsg) error {
	node := r.nodes[r.dataMap[key]]
	client, err := NewMqClient(fmt.Sprintf("%s:%d", node.Config.Name, node.Config.Port), 10*time.Second)
	if err != nil {
		errorMsg := fmt.Sprintf("Unable connect to node %s:%d\n", node.Config.Name, node.Config.Port)
		Logging(errorMsg, "ERROR")
		return errors.New(errorMsg)
	}

	err = client.CallDirect("GetItem", key, result)
	if err != nil {
		errorMsg := fmt.Sprintf("Unable to get data from node : %s", err.Error())
		return errors.New(errorMsg)
	}

	return nil
}

func (r *MqRPC) GetTable(key MqMsg, result *MqMsg) error {
	table := key.Key
	splitOwner := strings.Split(key.Value.(string), "|")
	filterOwner := splitOwner[1]
	ActiveUser := splitOwner[0]
	//fmt.Println("Owner: ", owner)
	//fmt.Println("Table: ", table)
	var tableContent []Table
	for k, v := range r.items {
		splitKey := strings.Split(k, "|")
		tableOwner := splitKey[0]
		tableName := splitKey[1]
		if tableName == table {
			row := Table{}
			row.Key = k
			row.Value = v.Value.(string)
			row.Owner = tableOwner
			if filterOwner == "" {
				if tableOwner == "public" || tableOwner == ActiveUser {
					tableContent = append(tableContent, row)
				}
			} else {
				if tableOwner == filterOwner {
					tableContent = append(tableContent, row)
				}
			}
		}
	}
	//table := Table{}
	buf, _ := Encode(tableContent)
	result.Value = buf.Bytes()
	return nil
}

func (r *MqRPC) GetWithBuildKey(key string, result *MqMsg) error {
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


func GetTableByKey(key string) string{
	tablePositionAtIndex := len(strings.Split(key, "|")) - 2
	tableName := strings.Split(key, "|")[tablePositionAtIndex]
	return tableName
}

func (r *MqRPC) setTableProperties(value MqMsg) {
	tableName := GetTableByKey(value.Key)
	table := NewTable(tableName,value.Owner)
	isTableExist := false
	for k,v := range r.tables{
			if k == tableName{
				*table = v
				isTableExist = true
				break
			}
	}
	item := make(map[string]interface{})
	if !isTableExist{
		item[value.Key] = value.Value
		table.Items = item
		r.tables[tableName] = *table
		message := fmt.Sprintf("Succesfull add new table %s and the properties ",tableName)
		Logging (message, "INFO")
	}else{
		table.Items[value.Key] = value.Value
		message := fmt.Sprintf("Succesfull add item, key->%s, value->%s, in table %s",value.Key,value.Value,tableName)
		Logging (message, "INFO")
	}
	setIndex(table)
	// fmt.Println(r.tables)
}

func setIndex(t *MqTable){
	getIndexByRole := func(value interface{}) string {return GetEmployeeRole(value)}
	t.RunIndex("employeerole",getIndexByRole)
}
