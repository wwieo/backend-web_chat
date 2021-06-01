package main

import (
	"flag"

	"github.com/chat"
)

var (
	port = flag.String("p", ":8080", "set port")
)

func init() {
	flag.Parse()
}

func main() {
	chat.Start(*port)
}

//let ws = new WebSocket("ws://localhost:8080/chat")
