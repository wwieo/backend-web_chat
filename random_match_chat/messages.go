package random_match_chat

import utils "backend-web_chat/utils"

type Message struct {
	IP     string "json:'ip'"
	ID     int64  "json:'id'"
	Body   string "json:'body'"
	Sender string "json:'sender'"
}

func NewMessage(body string, sender string) *Message {
	return &Message{
		IP:     utils.GetUserIP(),
		ID:     utils.GetRandomI64(),
		Body:   body,
		Sender: sender,
	}
}
