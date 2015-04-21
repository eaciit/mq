package monitor

import(
	"fmt"
	"os"
	"time"
	"path/filepath"
	"net/http"
	"html/template"
	"encoding/json"
	. "github.com/eaciit/mq/client"
	. "github.com/eaciit/mq/server"
)

type MqMonitor struct {
	port int
}

func (m *MqMonitor) Start() {
	client, e := NewMqClient("127.0.0.1:7890", time.Second * 10)
	handleError(e)

	http.Handle("/res/", http.StripPrefix("/res", http.FileServer(http.Dir(getView("assets")))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles(getView("views/index.gtpl"))
		t.Execute(w, nil)
	})

	http.HandleFunc("/test/nodes", func(w http.ResponseWriter, r *http.Request) {
	    for _, node := range getNodes(client) {
	    	fmt.Println(node.Config.Name)
	    }
	})

	http.ListenAndServe(fmt.Sprintf(":%d", m.port), nil)
}

func getView(view string) string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, "mqmonitor/web", view)
}

func getNodes(client *MqClient) []Node {
	msg, err := client.Call("AllNode", "")
	handleError(err)

	var nodes []Node
	if err := json.Unmarshal(msg.Value.([]byte), &nodes); err != nil {
		handleError(err)
    }

    return nodes
}

func handleError(e error) {
	if e != nil {
		panic(e.Error())
		os.Exit(100)
	}
}

func StartHTTP(port int) {
	monitor := new(MqMonitor);
	monitor.port = port
	monitor.Start()
}
