package random_match_chat

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type User struct {
	IP       string
	UserName string
	Conn     *websocket.Conn
	Global   *Chat
}

func (u *User) Read() {
	for {
		if _, message, err := u.Conn.ReadMessage(); err != nil {
			log.Println("Error on read messsage:", err.Error())
			break
		} else {
			u.Global.messages <- NewMessage(string(message), u.UserName)
		}
	}
	u.Global.leave <- u
}

func (u *User) Write(message *Message) {
	b, _ := json.Marshal(message)
	if err := u.Conn.WriteMessage(websocket.TextMessage, b); err != nil {
		log.Println("Error on writing message:", err.Error())
	}
}
