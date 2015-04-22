package helper

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
)

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
