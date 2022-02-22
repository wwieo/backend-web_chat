package message

import "time"

type Message struct {
	IP     string    `json:"ip" bson:"ip"`
	ID     string    `json:"id" bson:"_id"`
	Body   string    `json:"body" bson:"body"`
	Sender string    `json:"sender" bson:"sender"`
	Time   time.Time `json:"time" bson:"time"`
}
