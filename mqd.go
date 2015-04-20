package main

import (
	"flag"
	"fmt"
	. "github.com/eaciit/mq/client"
	. "github.com/eaciit/mq/server"
	"strings"
	"time"
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

	time.Sleep(10)
	c, e := NewMqClient("127.0.0.1:7890", time.Second*10)
	if e != nil {
		fmt.Println("Error: ", e.Error())
		return
	}
	status := ""
	t0 := time.Now()
	for status == "" {
		s, e := c.CallString("GetLog", t0)
		if e != nil {
			fmt.Println("Error: ", e.Error())
			status = "exit"
		} else {
			if s != "" {
				fmt.Println(s)
			}
			if strings.Contains(s, "EXIT") {
				status = "exit"
			}
		}
		t0 = time.Now()
	}
}
