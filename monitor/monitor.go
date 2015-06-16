package monitor

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/eaciit/mq/client"
	. "github.com/eaciit/mq/helper"
	. "github.com/eaciit/mq/msg"
	. "github.com/eaciit/mq/server"
)

const (
	ReconnectDelay   time.Duration = 3
	ConnectionTimout time.Duration = time.Second * 10
	ItemsLimit       int           = 50
	BaseView         string        = "monitor/web/"
	DevelopmentMode  bool          = true
)

var (
	ConnectionServerHost string
	Layout               *template.Template = GetTemplateView(BaseView + "views/*")
	clientInfo           *ClientInfo
)

type (
	MqMonitor struct {
		port int
	}

	FuncParam func() error
)

func (m *MqMonitor) Start() {
	var client *MqClient
	var err error

	clientInfo = &ClientInfo{IsLoggedIn: false}

	client, err = connect()
	Errorable(err)

	http.Handle("/res/", http.StripPrefix("/res", http.FileServer(http.Dir(GetView(BaseView+"assets")))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleIndex(w, r, client, err)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handleLogin(w, r, client, err)
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		handleLogout(w, r, client, err)
	})

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		handleUser(w, r, client, err)
	})

	http.HandleFunc("/console", func(w http.ResponseWriter, r *http.Request) {
		handleConsole(w, r, client, err)
	})

	http.HandleFunc("/data/nodes", func(w http.ResponseWriter, r *http.Request) {
		handleDataNodes(w, r, client, err)
	})

	http.HandleFunc("/data/items", func(w http.ResponseWriter, r *http.Request) {
		handleDataItems(w, r, client, err)
	})

	http.HandleFunc("/data/users", func(w http.ResponseWriter, r *http.Request) {
		handleDataUsers(w, r, client, err)
	})

	fmt.Printf("starting http at :%d, connecting to master %s\n", m.port, ConnectionServerHost)
	err = http.ListenAndServe(fmt.Sprintf(":%d", m.port), nil)
	Errorable(err, func() {
		if err != nil {
			os.Exit(0)
		}
	})
}

func handleIndex(w http.ResponseWriter, r *http.Request, client *MqClient, err error) {
	if !clientInfo.IsLoggedIn {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	executeTemplate(w, "index", nil)
}

func handleLogin(w http.ResponseWriter, r *http.Request, client *MqClient, err error) {
	if r.Method != "GET" {
		w.Header().Set("Content-type", "application/json")
		r.ParseForm()
	}

	if isServerAlive(w, r, client) == false {
		return
	}

	if r.Method == "GET" {
		if clientInfo.IsLoggedIn {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		executeTemplate(w, "login", nil)
	} else if r.Method == "POST" {
		if clientInfo.IsLoggedIn {
			PrintJSON(w, false, "", "already logged in")
			return
		}

		if success := rpcDo(w, client, func() error {
			msg, err := client.CallToLogin(MqMsg{
				Key:   r.FormValue("username"),
				Value: r.FormValue("password"),
			})

			*clientInfo = msg.Value.(ClientInfo)

			if err != nil {
				return err
			}

			if !clientInfo.IsLoggedIn {
				return errors.New("Username/password is not match")
			}

			return err
		}); !success {
			return
		}

		PrintJSON(w, true, clientInfo, "")
		return
	}
}

func handleLogout(w http.ResponseWriter, r *http.Request, client *MqClient, err error) {
	clientInfo = &ClientInfo{IsLoggedIn: false}

	if r.Method == "GET" {
		w.Header().Set("Content-type", "application/json")

		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	PrintJSON(w, false, "", "you are not logged in. login first")
}

func handleUser(w http.ResponseWriter, r *http.Request, client *MqClient, err error) {
	if !clientInfo.IsLoggedIn {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	executeTemplate(w, "user", nil)
}

func handleConsole(w http.ResponseWriter, r *http.Request, client *MqClient, err error) {
	if r.Method != "GET" {
		w.Header().Set("Content-type", "application/json")
		r.ParseForm()
	}

	if isServerAlive(w, r, client) == false {
		return
	}

	if r.Method == "GET" {
		if !clientInfo.IsLoggedIn {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		executeTemplate(w, "console", nil)
	} else if r.Method == "POST" {
		if !clientInfo.IsLoggedIn {
			PrintJSON(w, false, "", "you are not logged in. login first")
			return
		}

		mode := strings.ToLower(r.FormValue("mode"))
		key := strings.ToLower(r.FormValue("key"))
		value := strings.ToLower(r.FormValue("value"))
		owner := strings.ToLower(r.FormValue("owner"))
		table := strings.ToLower(r.FormValue("table"))
		duration := strings.ToLower(r.FormValue("duration"))
		permission := strings.ToLower(r.FormValue("permission"))
		keyParsed := ""

		_ = duration
		_ = permission

		if owner == "" {
			owner = "public"
		}

		if owner != "" {
			keyParsed = fmt.Sprintf("%s", owner)
		}

		if table != "" {
			keyParsed = fmt.Sprintf("%s|%s", keyParsed, table)
		}

		if key != "" {
			keyParsed = fmt.Sprintf("%s|%s", keyParsed, key)
		}

		if mode == "get" {
			rpcDo(w, client, func() error {
				msg, err := client.Call("Get", keyParsed)

				if err == nil {
					PrintJSON(w, true, msg.Value, "")
				}

				return err
			})

			return
		} else if mode == "set" {
			rpcDo(w, client, func() error {
				_, err := client.Call("Set", MqMsg{
					Key:   keyParsed,
					Value: value,
				})

				if err == nil {
					PrintJSON(w, true, keyParsed, "")
				}

				return err
			})

			return
		}

		PrintJSON(w, false, "", "Bad request")
		return
	}
}

func handleDataNodes(w http.ResponseWriter, r *http.Request, client *MqClient, err error) {
	w.Header().Set("Content-type", "application/json")
	r.ParseForm()

	if !clientInfo.IsLoggedIn {
		PrintJSON(w, false, "", "you are not logged in. login first")
		return
	}

	if isServerAlive(w, r, client) == false {
		return
	}

	dataSizeUnit, _ := strconv.ParseInt(r.FormValue("dataSizeUnit"), 10, 64)

	if r.Method == "GET" {
		var nodes []Node

		if success := rpcDo(w, client, func() error {
			return client.CallDecode("Nodes", "", &nodes)
		}); !success {
			return
		}

		searchKeyword := strings.ToLower(r.FormValue("search"))
		var resultGrid []map[string]interface{}

		for _, node := range nodes {
			dataNode := map[string]interface{}{
				"ConfigName":    node.Config.Name,
				"ConfigPort":    node.Config.Port,
				"ConfigRole":    node.Config.Role,
				"DataCount":     node.DataCount,
				"DataSize":      node.DataSize,
				"AllocatedSize": node.AllocatedSize / dataSizeUnit,
				"StartTime":     node.StartTime.Format("2006-01-02 15:04:05"),
				"Duration":      FormatDuration(time.Since(node.StartTime)),
			}

			isExist := (len(searchKeyword) == 0)
			for _, v := range dataNode {
				if strings.Contains(strings.ToLower(AsString(v)), searchKeyword) {
					isExist = true
					break
				}
			}

			if isExist {
				resultGrid = append(resultGrid, dataNode)
			}
		}

		seriesLimit, _ := strconv.Atoi(r.FormValue("seriesLimit"))
		seriesDelay, _ := strconv.Atoi(r.FormValue("seriesDelay"))
		nowTime := time.Now()
		resultChart := make([]map[string]interface{}, seriesLimit)

		for i := 0; i < seriesLimit; i++ {
			eachNowTime := nowTime.Add(time.Duration(-i*seriesDelay) * time.Second)
			eachNowTimeInt, err := strconv.ParseInt(fmt.Sprintf("%s", eachNowTime.Format("02150405")), 10, 32)

			if err != nil {
				continue
			}

			var totalHost int64 = 0
			var totalDataCount int64 = 0
			var totalDataSize int64 = 0
			var totalAllocatedSize int64 = 0

			for _, node := range nodes {
				eachNodeTimeInt, err := strconv.ParseInt(fmt.Sprintf("%s", node.StartTime.Format("02150405")), 10, 32)

				if err != nil {
					continue
				}

				if eachNodeTimeInt <= eachNowTimeInt {
					totalHost += 1
					totalDataCount += node.DataCount
					totalDataSize += node.DataSize
					totalAllocatedSize += (node.AllocatedSize / dataSizeUnit)
				}
			}

			resultChart[i] = map[string]interface{}{
				"Time":               eachNowTime.Format("15:04:05"),
				"TimeInt":            eachNowTimeInt,
				"TotalHost":          totalHost,
				"TotalDataCount":     totalDataCount,
				"TotalDataSize":      totalDataSize,
				"TotalAllocatedSize": totalAllocatedSize,
			}
		}

		result := map[string]interface{}{
			"grid":  resultGrid,
			"chart": resultChart,
		}

		PrintJSON(w, true, result, "")
		return
	}

	PrintJSON(w, false, "", "Bad Request")
}

func handleDataItems(w http.ResponseWriter, r *http.Request, client *MqClient, err error) {
	w.Header().Set("Content-type", "application/json")
	r.ParseForm()

	if !clientInfo.IsLoggedIn {
		PrintJSON(w, false, "", "you are not logged in. login first")
		return
	}

	if isServerAlive(w, r, client) == false {
		return
	}

	if r.Method == "GET" {
		var items map[string]MqMsg

		if success := rpcDo(w, client, func() error {
			return client.CallDecode("Items", "", &items)
		}); !success {
			return
		}

		searchKeyword := strings.ToLower(r.FormValue("search"))
		var resultGrid []map[string]interface{}

		i := 0
		for _, v := range items {
			if !(i < ItemsLimit) {
				break
			}

			dataNode := map[string]interface{}{
				"Key":        v.Key,
				"Value":      v.Value.(string),
				"Created":    v.Created.Format("2006-01-02 15:04:05"),
				"LastAccess": v.LastAccess.Format("2006-01-02 15:04:05"),
				"Expiry":     FormatDuration(v.Expiry),
			}

			isExist := (len(searchKeyword) == 0)
			for _, w := range dataNode {
				if strings.Contains(strings.ToLower(AsString(w)), searchKeyword) {
					isExist = true
					break
				}
			}

			if isExist {
				resultGrid = append(resultGrid, dataNode)
			}

			i += 1
		}

		result := map[string]interface{}{
			"grid": resultGrid,
		}

		PrintJSON(w, true, result, "")
		return
	}

	PrintJSON(w, true, make([]interface{}, 0), "")
}

func handleDataUsers(w http.ResponseWriter, r *http.Request, client *MqClient, err error) {
	w.Header().Set("Content-type", "application/json")
	r.ParseForm()

	if !clientInfo.IsLoggedIn {
		PrintJSON(w, false, "", "you are not logged in. login first")
		return
	}

	if isServerAlive(w, r, client) == false {
		return
	}

	if r.Method == "GET" {
		var users []MqUser

		if success := rpcDo(w, client, func() error {
			return client.CallDecode("Users", "", &users)
		}); !success {
			return
		}

		searchKeyword := strings.ToLower(r.FormValue("search"))
		var resultGrid []map[string]interface{}

		for _, v := range users {
			dataUser := map[string]interface{}{
				"UserName": v.UserName,
				"Role":     v.Role,
			}

			isExist := (len(searchKeyword) == 0)
			for _, w := range dataUser {
				if w == "Password" {
					continue
				}

				if strings.Contains(strings.ToLower(AsString(w)), searchKeyword) {
					isExist = true
					break
				}
			}

			if isExist {
				resultGrid = append(resultGrid, dataUser)
			}
		}

		result := map[string]interface{}{
			"grid": resultGrid,
		}

		PrintJSON(w, true, result, "")
		return
	} else if r.Method == "POST" {
		oldUsername := strings.ToLower(r.FormValue("oldusername"))
		username := strings.ToLower(r.FormValue("username"))
		password := strings.ToLower(r.FormValue("password"))
		role := strings.ToLower(r.FormValue("role"))
		isEdit := strings.ToLower(r.FormValue("edit"))

		if isEdit == "true" {
			if success := rpcDo(w, client, func() error {
				_, e := client.Call("ChangePassword", MqMsg{
					Key:   oldUsername,
					Value: password,
				})

				return e
			}); !success {
				return
			}

			PrintJSON(w, true, make([]interface{}, 0), "")
			return
		}

		if success := rpcDo(w, client, func() error {
			_, e := client.Call("AddUser", MqMsg{
				Key:   fmt.Sprintf("%s|%s", username, role),
				Value: password,
			})

			return e
		}); !success {
			return
		}

		PrintJSON(w, true, make([]interface{}, 0), "")
		return
	} else if r.Method == "DELETE" {
		if clientInfo.Role == "admin" {
			username := strings.ToLower(r.FormValue("username"))

			if success := rpcDo(w, client, func() error {
				_, e := client.Call("DeleteUser", MqMsg{
					Key:   username,
					Value: username,
				})

				return e
			}); !success {
				return
			}

			PrintJSON(w, true, make([]interface{}, 0), "")
		} else {
			PrintJSON(w, false, "", "You don't have permission to delete")
		}
		return
	}

	PrintJSON(w, true, make([]interface{}, 0), "")
}

func connect() (*MqClient, error) {
	return NewMqClient(ConnectionServerHost, ConnectionTimout)
}

func executeTemplate(w http.ResponseWriter, page string, data interface{}) {
	if DevelopmentMode {
		Layout = GetTemplateView(BaseView + "views/*")
	}

	Layout.ExecuteTemplate(w, page, data)
}

func rpcDo(w http.ResponseWriter, client *MqClient, fn FuncParam) bool {
	if client == nil {
		PrintJSON(w, false, "", "connection is shut down")
		return false
	}

	if err := fn(); err != nil {
		PrintJSON(w, false, "", err.Error())
		return false
	}

	return true
}

func isServerAlive(w http.ResponseWriter, r *http.Request, client *MqClient) bool {
	if isServerAlive := r.FormValue("isServerAlive"); isServerAlive == "false" {
		var err error
		client, err = connect()

		if err != nil {
			fmt.Println(err.Error())
			PrintJSON(w, false, "", "connection is shut down")
			return false
		}
	}

	return true
}

func StartHTTP(serverHost string, port int) {
	ConnectionServerHost = serverHost

	monitor := new(MqMonitor)
	monitor.port = port
	monitor.Start()
}
