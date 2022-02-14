package main

import (
	"flag"

	random_match_chat "backend-web_chat/random_match_chat"
)

var (
	port = flag.String("p", ":8080", "set port")
)

func init() {
	flag.Parse()
}

func main() {
	random_match_chat.Start(*port)
}

//let ws = new WebSocket("ws://localhost:8080/chat")
