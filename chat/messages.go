package chat

import (
	utils "backend-web_chat/utils"
	"encoding/json"
)

type Message struct {
	IP     string `json:"IP"`
	ID     string `json:"ID"`
	Body   string `json:"Body"`
	Sender string `json:"Sender"`
}

func NewMessage(body string, sender string) []byte {
	result, _ := json.Marshal(
		&Message{
			IP:     utils.GetUserIP(),
			ID:     utils.GetMsgID(),
			Body:   body,
			Sender: sender,
		})
	return result
}
