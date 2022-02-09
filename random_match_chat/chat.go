package random_match_chat

import (
	"fmt"
	"log"
)

func (chat *Chat) Run() {
	for {
		select {
		case user := <-chat.join:
			chat.add(user)
		case message := <-chat.messages:
			chat.broadcast(message)
		case user := <-chat.leave:
			chat.disconnect(user)
		}
	}
}

func (chat *Chat) add(user *User) {
	if _, ok := chat.users[user.UserName]; !ok {
		chat.users[user.UserName] = user

		body := fmt.Sprintf("%s join the chat", user.UserName)
		chat.broadcast(NewMessage(body, "Server"))
	}
}

func (chat *Chat) broadcast(message *Message) {
	log.Printf("%s: %s\n", message.Sender, message.Body)
	for _, user := range chat.users {
		user.Write(message)
	}
}

func (chat *Chat) disconnect(user *User) {
	if _, ok := chat.users[user.UserName]; ok {
		defer user.Conn.Close()
		body := fmt.Sprintf("%s left the chat", user.UserName)
		delete(chat.users, user.UserName)
		chat.broadcast(NewMessage(body, "Server"))
	}
}
