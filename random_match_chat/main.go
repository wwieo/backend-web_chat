package random_match_chat

import (
	utils "backend-web_chat/utils"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func Start(port string) {
	ginDefault := gin.Default()
	ginDefault.GET("/", func(c *gin.Context) {
		c.Data(200, "text/plain", []byte("Hello webchat world!"))
	})

	melodyNew := melody.New()
	ginDefault.GET("/chat", func(c *gin.Context) {
		melodyNew.HandleRequest(c.Writer, c.Request)
	})

	melodyNew.HandleMessage(func(session *melody.Session, msg []byte) {
		melodyNew.Broadcast(NewMessage(string(msg), utils.GetUsername(session)))
	})

	melodyNew.HandleConnect(func(session *melody.Session) {
		melodyNew.Broadcast(NewMessage(utils.GetUsername(session)+" join the chat", "Server"))
	})

	melodyNew.HandleClose(func(session *melody.Session, i int, s string) error {
		melodyNew.Broadcast(NewMessage(utils.GetUsername(session)+" left the room", "Server"))
		return nil
	})
	ginDefault.Run(port)
}
