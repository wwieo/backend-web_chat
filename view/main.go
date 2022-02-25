package view

import (
	controllers "backend-web_chat/controllers"
	mdDB "backend-web_chat/model/database"
	"fmt"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func Start() {
	utilsController := controllers.NewUtilsController()
	config := utilsController.GetConfig()

	ginDefault := gin.Default()
	ginDefault.GET("/", func(c *gin.Context) {
		c.Data(200, "text/plain", []byte("Hello webchat world!"))
	})

	socketTool := mdDB.RedisTool{
		Melody:      melody.New(),
		RedisClient: controllers.SetRedisClient(),
	}
	chatController := controllers.NewChatController(&socketTool)

	mongoClient := controllers.SetMongoClient()
	mongoDB := config.GetString("mongo.database")
	mongoDBCol := config.GetString("mongo.collection")
	mongoTool := mdDB.MongoTool{
		MongoClient: mongoClient,
		Database:    mongoClient.Database(mongoDB),
		CollName:    mongoClient.Database(mongoDB).Collection(mongoDBCol),
	}

	ginDefault.GET("/chat", controllers.HandleRequest(&socketTool))

	chatController.HandleMessage(&socketTool, &mongoTool)
	chatController.HandleConnect(&socketTool, &mongoTool)
	chatController.HandleClose(&socketTool, &mongoTool)

	ginDefault.Run(fmt.Sprintf(":%d", config.GetInt("application.port")))
}
