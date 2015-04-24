package monitor

import (
	"fmt"
	. "github.com/eaciit/mq/client"
	. "github.com/eaciit/mq/helper"
	. "github.com/eaciit/mq/msg"
	. "github.com/eaciit/mq/server"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

const (
	ReconnectDelay       time.Duration = 3
	ConnectionServerHost string        = "127.0.0.1:7890"
	ConnectionTimout     time.Duration = time.Second * 10
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

	client, _ = connect()

	http.Handle("/res/", http.StripPrefix("/res", http.FileServer(http.Dir(GetView("mqmonitor/web/assets")))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles(GetView("mqmonitor/web/views/index.gtpl"))
		t.Execute(w, nil)
	})

	http.HandleFunc("/data/nodes", func(w http.ResponseWriter, r *http.Request) {
		dataNodes(w, r, client, err)
	})

	http.HandleFunc("/data/items", func(w http.ResponseWriter, r *http.Request) {
		dataItems(w, r, client, err)
	})

	http.ListenAndServe(fmt.Sprintf(":%d", m.port), nil)
}

func dataNodes(w http.ResponseWriter, r *http.Request, client *MqClient, err error) {
	w.Header().Set("Content-type", "application/json")
	r.ParseForm()

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

		resultGrid := make([]map[string]interface{}, len(nodes))

		for i, node := range nodes {
			resultGrid[i] = map[string]interface{}{
				"ConfigName":    node.Config.Name,
				"ConfigPort":    node.Config.Port,
				"ConfigRole":    node.Config.Role,
				"DataCount":     node.DataCount,
				"DataSize":      node.DataSize,
				"AllocatedSize": node.AllocatedSize / dataSizeUnit,
				"StartTime":     node.StartTime.Format("2006-01-02 15:04:05"),
				"Duration":      FormatDuration(time.Since(node.StartTime)),
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

func dataItems(w http.ResponseWriter, r *http.Request, client *MqClient, err error) {
	w.Header().Set("Content-type", "application/json")
	r.ParseForm()

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

		resultGrid := make([]map[string]interface{}, len(items))

		var i int = 0
		for _, v := range items {
			resultGrid[i] = map[string]interface{}{
				"Key":        v.Key,
				"Value":      v.Value.(string),
				"Created":    v.Created.Format("2006-01-02 15:04:05"),
				"LastAccess": v.LastAccess.Format("2006-01-02 15:04:05"),
				"Expiry":     FormatDuration(v.Expiry),
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

func connect() (*MqClient, error) {
	return NewMqClient(ConnectionServerHost, ConnectionTimout)
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

func StartHTTP(port int) {
	monitor := new(MqMonitor)
	monitor.port = port
	monitor.Start()
}
