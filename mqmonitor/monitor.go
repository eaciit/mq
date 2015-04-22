package monitor

import (
	"fmt"
	. "github.com/eaciit/mq/client"
	. "github.com/eaciit/mq/helper"
	. "github.com/eaciit/mq/server"
	"html/template"
	"net/http"
	"os"
	"time"
)

type MqMonitor struct {
	port int
}

func (m *MqMonitor) Start() {
	client, err := NewMqClient("127.0.0.1:7890", time.Second*10)
	handleError(err)

	http.Handle("/res/", http.StripPrefix("/res", http.FileServer(http.Dir(GetView("mqmonitor/web/assets")))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles(GetView("mqmonitor/web/views/index.gtpl"))
		t.Execute(w, nil)
	})

	http.HandleFunc("/data/nodes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		if r.Method == "GET" {
			var nodes []Node
			client.CallDecode("Nodes", "", &nodes)

			result := make([]map[string]interface{}, len(nodes))

			for i, node := range nodes {
				result[i] = map[string]interface{}{
					"ConfigName": node.Config.Name,
					"ConfigPort": node.Config.Port,
					"ConfigRole": node.Config.Role,
					"DataCount":  node.DataCount,
					"DataSize":   node.DataSize,
					"StartTime":  node.StartTime.Format("2006-01-02 15:04"),
					"Duration":   FormatDuration(node.StartTime),
				}
			}

			PrintJSON(w, true, result, "")
			return
		}

		PrintJSON(w, false, "", "Bad Request")
	})

	http.ListenAndServe(fmt.Sprintf(":%d", m.port), nil)
}

func handleError(e error) {
	if e != nil {
		panic(e.Error())
		os.Exit(100)
	}
}

func StartHTTP(port int) {
	monitor := new(MqMonitor)
	monitor.port = port
	monitor.Start()
}
