package main

import (
	"flag"
	"fmt"
	. "github.com/eaciit/mq/client"
	. "github.com/eaciit/mq/server"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main() {
	var e error

	runtime.GOMAXPROCS(runtime.NumCPU())
	portFlag := flag.Int("port", 7890, "Port of RCP call. Default is 7890")
	hostFlag := flag.String("master", "", "Master host. Default is localhost:7890")
	flag.Parse()

	hostName := "127.0.0.1"
	hostPort := 7890
	if *hostFlag != "" {
		hostParts := strings.Split(*hostFlag, ":")
		if hostParts[0] == "" {
			hostName = "127.0.0.1"
		}
		if len(hostParts) > 1 {
			if !(hostParts[1] == "" || hostParts[1] == "0") {
				hostPort, e = strconv.Atoi(hostParts[1])
				if e != nil {
					panic("Invalid master listener address")
				}
			}
		}
	}

	startStatus := make(chan string)
	fmt.Printf("Starting MQ server at port %d \n", *portFlag)
	go func() {
		e = StartMQServer("127.0.0.1", *portFlag)
		if e != nil {
			//panic("Unable to start server: " + e.Error())
			startStatus <- fmt.Sprintf("\nUnable to start service : %s \n", e.Error())
			return
		}
	}()

	time.Sleep(5 * time.Second)
	currentListenerAddress := "127.0.0.1:" + strconv.Itoa(*portFlag)
	c, e := NewMqClient(currentListenerAddress, time.Second*10)
	if e != nil {
		fmt.Println("Error: ", e.Error())
		return
	}
	defer c.Close()

	if *hostFlag != "" {
		cfg, e := c.Call("GetConfig", "")
		if e != nil {
			fmt.Println("Unable to get config : " + e.Error())
			return
		}
		s, e := NewMqClient(fmt.Sprintf("%s:%d", hostName, hostPort), time.Second*10)
		if e != nil {
			fmt.Printf("Unable to connect to master server %s:%d : %s \n", hostName, hostPort, e.Error())
			return
		}
		defer s.Close()
		_, e = s.Call("AddNode", cfg.Value.(ServerConfig))
		if e != nil {
			fmt.Printf("Unable to set as node : %s", e.Error())
			return
		}
		//-- c.Call("SetHost",&ServerConfig{})
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

		if *hostFlag != "" {
			//this is slave
			s, _ = c.CallString("CheckHealthMaster", fmt.Sprintf("%s:%d", hostName, hostPort))
			if s == "KILL" {
				status = "exit"
			}

		} else {
			//this is master
			c.CallString("CheckHealthSlaves", "")
		}

		t0 = time.Now()
		time.Sleep(1 * time.Second)
	}

}
