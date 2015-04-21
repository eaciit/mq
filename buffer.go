package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type data struct {
	V1 int
	V2 int
	S1 string
}

func main() {
	var e error
	datas := []data{}
	fmt.Println("Preparing data")
	fmt.Println("================")
	for i := 1; i <= 1000000; i++ {
		datas = append(datas, data{i, i * 2, fmt.Sprintf("Ini data ke %d", i)})
	}
	/*
		for _, d := range datas {
			fmt.Printf("%s\t\t | %d\t | %d \n", d.S1, d.V1, d.V2)
		}
	*/
	//--- change to byte[]
	t0 := time.Now()
	buf, e := Encode(datas)
	handleError(e, 100)
	fmt.Printf("\nEncoding success\n===============\n")
	fmt.Printf("Size(MB) : %2d \n", (buf.Len() / 1024 / 1024))

	fmt.Println("\nDecoding")
	fmt.Println("=====================")
	results := []data{}
	e = Decode(buf, &results)
	handleError(e, 100)
	fmt.Println("Time: ", time.Since(t0))

	t0 = time.Now()
	bs, e := json.Marshal(datas)
	handleError(e, 100)
	fmt.Printf("\nMarshalling success\n===============\n")
	fmt.Printf("Size(MB) : %2d \n", (len(bs) / 1024 / 1024))

	fmt.Println("\nUnmarshall")
	fmt.Println("=====================")
	results = []data{}
	e = json.Unmarshal(bs, &results)
	handleError(e, 100)
	/*
		for _, d := range results {
			fmt.Printf("%s\t\t | %d\t | %d \n", d.S1, d.V1, d.V2)
		}
	*/
	fmt.Println("Time: ", time.Since(t0))
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

func Decode(buf *bytes.Buffer, result interface{}) error {
	dec := gob.NewDecoder(buf)
	e := dec.Decode(result)
	return e
}

func handleError(e error, exitCode int) {
	if e != nil {
		fmt.Println("Error :", e.Error())
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}
}
