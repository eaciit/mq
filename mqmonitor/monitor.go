package monitor

import (
	"fmt"
	. "github.com/eaciit/mq/client"
	. "github.com/eaciit/mq/helper"
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
		w.Header().Set("Content-type", "application/json")
		r.ParseForm()

		if isServerAlive := r.FormValue("isServerAlive"); isServerAlive == "false" {
			client, err = connect()

			if err != nil {
				fmt.Println(err.Error())
				PrintJSON(w, false, "", "connection is shut down")
				return
			}
		}

		var dataSizeUnit int64 = 1024 // kb

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
					"ConfigName": node.Config.Name,
					"ConfigPort": node.Config.Port,
					"ConfigRole": node.Config.Role,
					"DataCount":  node.DataCount,
					"DataSize":   node.DataSize / dataSizeUnit,
					"StartTime":  node.StartTime.Format("2006-01-02 15:04:05"),
					"Duration":   FormatDuration(node.StartTime),
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

				totalHost := 0
				totalDataCount := 0
				totalDataSize := 0

				for _, node := range nodes {
					eachNodeTimeInt, err := strconv.ParseInt(fmt.Sprintf("%s", node.StartTime.Format("02150405")), 10, 32)

					if err != nil {
						continue
					}

					if eachNodeTimeInt <= eachNowTimeInt {
						totalHost += 1
						totalDataCount += int(node.DataCount)
						totalDataSize += int(node.DataSize / dataSizeUnit)
					}
				}

				resultChart[i] = map[string]interface{}{
					"Time":           eachNowTime.Format("15:04:05"),
					"TimeInt":        eachNowTimeInt,
					"TotalHost":      totalHost,
					"TotalDataCount": totalDataCount,
					"TotalDataSize":  totalDataSize,
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
	})

	http.ListenAndServe(fmt.Sprintf(":%d", m.port), nil)
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

func StartHTTP(port int) {
	monitor := new(MqMonitor)
	monitor.port = port
	monitor.Start()
}
