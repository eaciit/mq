package main

import (
	"flag"
	"fmt"
	. "github.com/eaciit/mq/server"
)

func main() {
	var e error
	portFlag := flag.Int("port", 7890, "Port of RCP call. Default is 7890")
	flag.Parse()

	startStatus := make(chan string)
	fmt.Printf("Starting MQ server at port %d \n", *portFlag)
	go func() {
		e = StartMQServer("", *portFlag, startStatus)
		if e != nil {
			//panic("Unable to start server: " + e.Error())
			startStatus <- fmt.Sprintf("\nUnable to start service: %s \n", e.Error())
			return
		}
	}()

	status := ""
	for status == "" {
		select {
		case status = <-startStatus:
			fmt.Println(status)
		}
	}
}
