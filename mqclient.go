package main

import (
	"fmt"
	. "github.com/eaciit/mq/client"
	"time"
)

func main() {
	var e error
	c, e := NewMqClient("127.0.0.1:7890", time.Second*10)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	fmt.Println("Connected to RPC Server")
	for i := 0; i < 5; i++ {
		result, e := c.CallString("Info", "")
		if e != nil {
			fmt.Println("Error call RPC Method Info :", e.Error())
			return
		}
		fmt.Printf("Call Result: %v \n", result)
	}
}
