package routes

import (
	"github.com/ashikkabeer/messaging-api/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
    router := gin.Default()

    router.GET("/messages", handlers.RetrieveHistory)
    router.POST("/messages", handlers.SendMessage)
    router.PATCH("/messages/:message_id/read", handlers.MarkAsRead)

    return router
}
