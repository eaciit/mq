package main

import (
	"fmt"
	"net"
	//"bufio"
	"strings"
	//"bytes"
	. "github.com/eaciit/mq/client"
	. "github.com/eaciit/mq/server"
	"os"
	"regexp"
	"strconv"
	"time"
	//"math"
)

func parseSetCommand(command string) (string, string) {
	match, _ := regexp.MatchString("set()", command)

	if match == true {
		splitSet := strings.Split("set(asdasd,sds:1234)", "set(")[1]
		data := strings.TrimRight(splitSet, ")")
		//datas := strings.Split(data, ",")
		//s := make([]string, 3)

		//d := append(d, datas[0])
		//for i := 1; i < len(datas-1); i++ {
		//	d := append(d, datas[i])
		//}
		return "set", data

	} else {
		return "set", "" //nil, make([]string, 0)
	}
}

func main() {
	fmt.Println(setCommand("set(asdasd,sds:1234)"))
	fmt.Println(setCommand("set(asdasd,sds:1234)"))

	/*
		#1 Create Initial data want to transmit to node
		#2 Check load (with count data amount for each node)
		#3 take lowest count data node
		#4 push data to the lowest count data node
		#5 (every pushing always check another node and do #3) until data count <= 0

	*/
	//nodeloadcheck()
	//overload()

	/* input  := ""
	    fmt.Scanf("%v", &input)
	    if(input != ""	) {
	    	addressandport := "127.0.0.1:7890"
			fmt.Println("Address and Port: " + addressandport )
		    addx := strings.Split(addressandport, ":")
		    transmit(addx[0],addx[1])

	    } */

}

func pushtonode(node Node) {

}

func nodeloadcheck() {
	var e error
	c, e := NewMqClient("127.0.0.1:7890", time.Second*10)
	handleError(e)
	s, e := c.CallString("Ping", "")
	handleError(e)
	fmt.Printf(s)

	var nodes []Node
	err := c.CallDecode("Nodes", "", &nodes)
	handleError(err)

	var adddata int64
	adddata = 102

	indx := make([]int64, len(nodes)) //[] int
	for e := range nodes {
		//datacount := e.DataCount
		//fmt.Println( nodes[e].Config.Name)
		//fmt.Println( nodes[e].DataCount)
		//fmt.Println( nodes[e].(Node).DataSize)
		indx[e] = nodes[e].DataCount + adddata
		fmt.Println(nodes[e].DataCount + adddata)
		adddata = 10 + adddata
	}
	fmt.Println(" ========== Get Max ========== ")
	fmt.Println(" INDX ")
	fmt.Println(indx)
	fmt.Println(" INDX  ")
	tt := getMax(indx)
	fmt.Println("Max Data Done ", tt)

}

func getMax(ys []int64) int64 {
	var output int64
	output = 0
	fmt.Println(strconv.Itoa(len(ys)))

	for i := 0; i < len(ys); i++ {
		if i == 0 {
			output = ys[i]
			fmt.Println("First ", output)
		} else {
			if output <= ys[i] {
				fmt.Println("> 1  ", output)
				output = ys[i]
			}
		}
	}
	return output
}

func handleError(e error) {
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(100)
	}
}

func overload() {
	// do distrbution to lowest load node, after node load check
}

func transmit(address string, port string) {
	// send data to node
	fmt.Println("Dialing... ")
	conn, error := net.Dial("tcp", address+":"+port)
	//fmt.Println(conn)
	//fmt.Println(error)
	fmt.Println("Done... ")

	if error != nil {
		fmt.Println("handle error")
	}
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	fmt.Println("==============================")
	//bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Listening... ")
	ln, err := net.Listen("tcp", address+":"+port)
	fmt.Println("Listening Done ")

	if err != nil {
		// handle error
	}
	for {
		_, err := ln.Accept()
		//fmt.Println("Listening => " + ln)

		if err != nil {
			// handle error
		}
		//go handleConnection(conn)
	}

}
