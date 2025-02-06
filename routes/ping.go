package routes

import (
	"health/controllers"

	"github.com/gin-gonic/gin"
)

func PingRoute(router *gin.RouterGroup) {
	auth := router.Group("/ping")
	{
		auth.GET("", controllers.Ping)
	}
}
