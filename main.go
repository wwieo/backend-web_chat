package main

import (
	view "backend-web_chat/view"
	"flag"
)

var (
	port = flag.String("p", ":8080", "set port")
)

func init() {
	flag.Parse()
}

func main() {
	view.Start(*port)
}

//let ws = new WebSocket("ws://localhost:8080/chat")
