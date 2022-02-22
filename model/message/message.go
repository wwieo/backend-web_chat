package message

type Message struct {
	IP     string `json:"IP" bson:"ip"`
	ID     string `json:"ID" bson:"_id"`
	Body   string `json:"Body" bson:"body"`
	Sender string `json:"Sender" bson:"sender"`
}
