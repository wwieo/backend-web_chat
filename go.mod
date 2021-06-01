module chat

go 1.16

replace github.com/utils => ./utils

replace github.com/chat => ./chat

require (
	github.com/chat v0.0.1 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/utils v0.0.1 // indirect
)
