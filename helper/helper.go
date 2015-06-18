package helper

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Pair struct {
	First, Second interface{}
}

func Decode(bytesData []byte, result interface{}) error {
	buf := bytes.NewBuffer(bytesData)
	dec := gob.NewDecoder(buf)
	e := dec.Decode(result)
	return e
}

func Encode(obj interface{}) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	gw := gob.NewEncoder(buf)
	err := gw.Encode(obj)
	if err != nil {
		return buf, err
	}
	return buf, nil
}

func PrintJSON(w http.ResponseWriter, success bool, data interface{}, message string) {
	w.Header().Set("Content-type", "application/json")

	result, err := json.Marshal(map[string]interface{}{
		"success": success,
		"data":    data,
		"message": message,
	})

	if err == nil {
		fmt.Fprintf(w, "%s\n", result)
	} else {
		result, _ := json.Marshal(map[string]interface{}{
			"success": false,
			"data":    nil,
			"message": err.Error(),
		})

		fmt.Fprintf(w, "%s\n", result)
	}
}

func GetTemplateView(path string) *template.Template {
	return template.Must(template.ParseGlob(GetView(path)))
}

func GetView(view string) string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, view)
}

func FormatDuration(duration time.Duration) string {
	hours := int(math.Floor(duration.Hours()))
	minutes := int(math.Floor(math.Mod(duration.Minutes(), 60)))
	seconds := int(math.Floor(math.Mod(math.Mod(duration.Seconds(), 3600), 60)))
	return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
}

func Errorable(err error, callbacks ...func()) {
	if err != nil {
		fmt.Printf("Error %s\n", err.Error())
	}

	if len(callbacks) > 0 {
		callbacks[0]()
	}
}

func AsString(val interface{}) string {
	return fmt.Sprintf("%v", val)
}

func FloatToString(val interface{}) string {
	return fmt.Sprintf("%.2f", val)
}
