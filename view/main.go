package view

import (
	controllers "backend-web_chat/controllers"
	mdDB "backend-web_chat/model/database"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func Start(port string) {

	ginDefault := gin.Default()
	ginDefault.GET("/", func(c *gin.Context) {
		c.Data(200, "text/plain", []byte("Hello webchat world!"))
	})

	socketTool := mdDB.RedisTool{
		Melody:      melody.New(),
		RedisClient: controllers.SetRedisClient(),
	}
	// mongoClient := controllers.SetMongoClient()
	// mongoTool := mdDB.MongoTool{
	// 	MongoClient: mongoClient,
	// 	Database:    mongoClient.Database("ChatSystem"),
	// 	CollName:    mongoClient.Database("ChatSystem").Collection("Messages"),
	// }

	ginDefault.GET("/chat", controllers.HandleRequest(&socketTool))

	controllers.HandleMessage(&socketTool)
	controllers.HandleConnect(&socketTool)
	controllers.HandleClose(&socketTool)

	ginDefault.Run(port)
}
