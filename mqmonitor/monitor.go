package monitor

import(
	"fmt"
	"os"
	"path/filepath"
	"net/http"
	"html/template"
)

type Monitor struct {
	port int
}

func (m *Monitor) Start() {
	http.Handle("/res/", http.StripPrefix("/res", http.FileServer(http.Dir(getView("assets")))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles(getView("views/index.gtpl"))
		t.Execute(w, nil)
	})

	http.ListenAndServe(fmt.Sprintf(":%d", m.port), nil)
}

func getView(view string) string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, "mqmonitor/web", view)
}

func handleError(e error) {
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(100)
	}
}

func StartHTTP(port int) {
	monitor := new(Monitor);
	monitor.port = port
	monitor.Start()
}