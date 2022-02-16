package chat

import (
	match "backend-web_chat/match"
	utils "backend-web_chat/utils"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gopkg.in/olahol/melody.v1"
)

type Chat interface {
	HandleRequest()
	HandleMessage()
	HandleConnect()
	HandleClose()
}

type SocketTool struct {
	melody      *melody.Melody
	redisClient *redis.Client
}

//when user request the route
func (socketTool *SocketTool) HandleRequest() func(c *gin.Context) {
	return func(c *gin.Context) {
		socketTool.melody.HandleRequest(c.Writer, c.Request)
	}
}

//when user send the message
func (socketTool *SocketTool) HandleMessage() {
	melodyNew := socketTool.melody
	redisClient := socketTool.redisClient

	melodyNew.HandleMessage(func(session *melody.Session, msg []byte) {
		id := match.GetSessionID(session)
		chatTo, _ := redisClient.Get(context.TODO(), id).Result()
		msg = NewMessage(string(msg), utils.GetUsername(session))
		//socket chat filtered with sessionID and only allow 1 on 1 talk
		melodyNew.BroadcastFilter(msg, func(session *melody.Session) bool {
			compareID, _ := session.Get(match.KEY)
			return compareID == chatTo || compareID == id
		})
	})
}

//when user connect the chat
func (socketTool *SocketTool) HandleConnect() {
	melodyNew := socketTool.melody
	redisClient := socketTool.redisClient

	melodyNew.HandleConnect(func(session *melody.Session) {
		id := match.InitSession(session)
		//check if 2 people are waiting to chat
		if key, err := match.GetWaitFirstKey(redisClient); err == nil && key != "" {
			//create 1 on 1 chat
			match.CreateChat(id, key, redisClient)
			msg := NewMessage("Match success", "Server")
			melodyNew.BroadcastFilter(msg, func(session *melody.Session) bool {
				compareID, _ := session.Get(match.KEY)
				return compareID == id || compareID == key
			})
		} else {
			// if not then add to the queue and wait
			match.AddToWaitList(id, socketTool.redisClient)
			msg := NewMessage("Wait for matching...", "Server")
			melodyNew.BroadcastFilter(msg, func(session *melody.Session) bool {
				compareID, _ := session.Get(match.KEY)
				return compareID == id
			})
		}
	})
}

//when user disconnect
func (socketTool *SocketTool) HandleClose() {
	melodyNew := socketTool.melody
	redisClient := socketTool.redisClient

	melodyNew.HandleClose(func(session *melody.Session, i int, s string) error {
		id := match.GetSessionID(session)
		chatTo, _ := redisClient.Get(context.TODO(), id).Result()
		msg := NewMessage(utils.GetUsername(session)+" left the room", "Server")
		match.RemoveChat(id, chatTo, redisClient)
		return melodyNew.BroadcastFilter(msg, func(session *melody.Session) bool {
			compareID, _ := session.Get(match.KEY)
			return compareID == chatTo
		})
	})
}
