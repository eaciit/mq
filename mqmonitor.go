package main

import (
	"flag"
	"github.com/eaciit/mq/mqmonitor"
)

func main() {
	port := flag.Int("port", 1234, "Port of RCP call. Default is 1234")
	serverHost := flag.String("master", "127.0.0.1:7890", "Default master host")
	flag.Parse()

	monitor.StartHTTP(*serverHost, *port)
}
