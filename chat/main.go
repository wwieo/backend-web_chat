package chat

import (
	match "backend-web_chat/match"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func Start(port string) {

	ginDefault := gin.Default()
	ginDefault.GET("/", func(c *gin.Context) {
		c.Data(200, "text/plain", []byte("Hello webchat world!"))
	})

	socketTool := SocketTool{
		melody:      melody.New(),
		redisClient: match.SetRedisClient(),
	}

	ginDefault.GET("/chat", socketTool.HandleRequest())

	socketTool.HandleMessage()
	socketTool.HandleConnect()
	socketTool.HandleClose()

	ginDefault.Run(port)
}
