package main

import (
	"bufio"
	"fmt"
	. "github.com/eaciit/mq/client"
	. "github.com/eaciit/mq/msg"
	"os"
	"runtime"
	"strings"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var e error
	c, e := NewMqClient("127.0.0.1:7890", time.Second*10)
	handleError(e)
	fmt.Println("Connected to RPC Server")
	r := bufio.NewReader(os.Stdin)

	status := ""
	for status != "exit" {
		fmt.Print("> ")
		//fmt.Scanln(&command)
		line, _, _ := r.ReadLine()
		command := string(line)
		handleError(e)
		//fmt.Printf("Processing command: %v \n", command)
		lowerCommand := strings.ToLower(command)

		if lowerCommand == "exit" {
			status = "exit"
			c.Close()
		} else if lowerCommand == "kill" {
			c.CallString("Kill", "")
			status = "exit"
			c.Close()
		} else if lowerCommand == "ping" {
			s, e := c.CallString("Ping", "")
			handleError(e)
			fmt.Printf(s)
		} else if strings.HasPrefix(lowerCommand, "set") {
			//--- this to handle set command
			commandParts := strings.Split(command, " ")
			key := commandParts[1]
			value := strings.Join(commandParts[2:], " ")
			msg := MqMsg{Key: key, Value: value}
			_, e := c.Call("Set", msg)
			if e != nil {
				fmt.Println("Unable to store message: " + e.Error())
			}
		} else if strings.HasPrefix(lowerCommand, "get") {
			//--- this to handle set command
			commandParts := strings.Split(command, " ")
			key := commandParts[1]
			msg, e := c.Call("Get", key)
			if e != nil {
				fmt.Println("Unable to store message: " + e.Error())
			} else {
				fmt.Printf("Value: %v \n", msg.Value)
			}
		}
	}
}

func handleError(e error) {
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(100)
	}
}
