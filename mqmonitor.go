package main

import (
	"flag"
	"github.com/eaciit/mq/mqmonitor"
)

func main() {
	port := flag.Int("port", 1234, "-port=1234")
	flag.Parse()
	monitor.StartHTTP(*port)
}
