package chat

import (
	match "backend-web_chat/match"
	utils "backend-web_chat/utils"
	"context"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func Start(port string) {
	redisClient := match.SetRedisClient()

	ginDefault := gin.Default()
	ginDefault.GET("/", func(c *gin.Context) {
		c.Data(200, "text/plain", []byte("Hello webchat world!"))
	})

	melodyNew := melody.New()
	ginDefault.GET("/chat", func(c *gin.Context) {
		melodyNew.HandleRequest(c.Writer, c.Request)
	})

	//when user send message
	melodyNew.HandleMessage(func(session *melody.Session, msg []byte) {
		id := match.GetSessionID(session)
		chatTo, _ := redisClient.Get(context.TODO(), id).Result()
		msg = NewMessage(string(msg), utils.GetUsername(session))
		melodyNew.BroadcastFilter(msg, func(session *melody.Session) bool {
			compareID, _ := session.Get(match.KEY)
			return compareID == chatTo || compareID == id
		})
	})

	//when user connect
	melodyNew.HandleConnect(func(session *melody.Session) {
		id := match.InitSession(session)
		if key, err := match.GetWaitFirstKey(redisClient); err == nil && key != "" {
			match.CreateChat(id, key, redisClient)
			msg := NewMessage(utils.GetUsername(session)+" join the chat", "Server")
			melodyNew.BroadcastFilter(msg, func(session *melody.Session) bool {
				compareID, _ := session.Get(match.KEY)
				return compareID == id || compareID == key
			})
		} else {
			match.AddToWaitList(id, redisClient)
		}
	})

	//when user disconnect
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
	ginDefault.Run(port)
}
