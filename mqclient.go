package main

import (
	"bufio"
	"fmt"
	. "github.com/eaciit/mq/client"
	. "github.com/eaciit/mq/msg"
	. "github.com/eaciit/mq/server"
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

	msg := MqMsg{Key: "INFO", Value: "New Client Connected"}
	c.CallToLog("SetLog", msg)

	r := bufio.NewReader(os.Stdin)

	status := ""
	for status != "exit" {
		fmt.Print("> ")
		//fmt.Scanln(&command)
		line, _, _ := r.ReadLine()
		command := string(line)
		handleError(e)
		//fmt.Printf("Processing command: %v \n", command)
		stringsPart := strings.Split(command, " ")
		lowerCommand := strings.ToLower(stringsPart[0])

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
		} else if lowerCommand == "nodes" {
			results := []Node{}
			e := c.CallDecode("Nodes", "", &results)
			handleError(e)
			fmt.Printf("%v\n", results)
		} else if lowerCommand == "set" {
			//--- this to handle set command
			commandParts := strings.Split(command, " ")
			key := commandParts[1]
			value := strings.Join(commandParts[2:], " ")
			msg := MqMsg{Key: key, Value: value}
			_, e := c.Call("Set", msg)
			if e != nil {
				fmt.Println("Unable to store message: " + e.Error())
			}
		} else if lowerCommand == "get" {
			//--- this to handle get command
			commandParts := strings.Split(command, " ")
			key := commandParts[1]
			msg, e := c.Call("Get", key)
			if e != nil {
				fmt.Println("Unable to store message: " + e.Error())
			} else {
				fmt.Printf("Value: %v \n", msg.Value)
			}
		} else if lowerCommand == "getlog" {
			commandParts := strings.Split(command, " ")
			key := commandParts[1]
			value := strings.Join(commandParts[2:], " ")
			msg := MqMsg{Key: key, Value: value}
			s, e := c.CallString("GetLogData", msg)
			handleError(e)
			fmt.Println(s)
		} else if lowerCommand == "adduser" {
			//--- this to handle set command
			commandParts := strings.Split(command, " ")
			userName := commandParts[1]
			password := strings.Join(commandParts[2:], " ")
			msg := MqMsg{Key: userName, Value: password}
			_, e := c.Call("AddUser", msg)
			if e != nil {
				fmt.Println("Unable to store message: " + e.Error())
			}
		} else if lowerCommand == "getlistusers" {
			s, e := c.CallString("GetListUsers", "")
			handleError(e)
			fmt.Printf(s)
		} else {
			errorMsg := "Unable to find command " + command
			//c.CallToLog(errorMsg,"ERROR")
			msg := MqMsg{Key: "ERROR", Value: errorMsg}
			_, e := c.CallToLog("SetLog", msg)
			if e != nil {
				fmt.Println("Unable to store message: " + e.Error())
			}

			fmt.Println(errorMsg)
		}
	}
}

func handleError(e error) {
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(100)
	}
}
