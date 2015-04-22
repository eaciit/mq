package monitor

import (
	"fmt"
	. "github.com/eaciit/mq/client"
	. "github.com/eaciit/mq/helper"
	. "github.com/eaciit/mq/server"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type MqMonitor struct {
	port int
}

func (m *MqMonitor) Start() {
	client, e := NewMqClient("127.0.0.1:7890", time.Second*10)
	handleError(e)

	http.Handle("/res/", http.StripPrefix("/res", http.FileServer(http.Dir(getView("assets")))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles(getView("views/index.gtpl"))
		t.Execute(w, nil)
	})

	http.HandleFunc("/nodes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		if r.Method == "GET" {
			nodes := getNodes(client)

			PrintJSON(w, true, nodes, "")
			return
		}

		PrintJSON(w, false, "", "Bad Request")
	})

	http.ListenAndServe(fmt.Sprintf(":%d", m.port), nil)
}

func getView(view string) string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, "mqmonitor/web", view)
}

func getNodes(client *MqClient) []Node {
	var nodes []Node

	err := client.CallDecode("Nodes", "", &nodes)
	handleError(err)

	return nodes
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
