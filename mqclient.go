package main

import (
	"fmt"
	. "github.com/eaciit/mq/client"
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

	status := ""
	for status != "exit" {
		command := ""
		fmt.Print("> ")
		fmt.Scanln(&command)
		handleError(e)
		command = strings.ToLower(command)

		if command == "exit" {
			status = "exit"
		} else if command == "kill" {
			c.CallString("Kill", "")
			status = "exit"
		} else if command == "ping" {
			s, e := c.CallString("Ping", "")
			handleError(e)
			fmt.Printf(s)
		}
	}
}

func handleError(e error) {
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(100)
	}
}
