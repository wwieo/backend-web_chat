package chat

import "github.com/gorilla/websocket"

type User struct {
	UserName string
	Conn     *websocket.Conn
	Global   *Chat
}
