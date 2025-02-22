package routes

import (
	"health/controllers"
	"health/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.RouterGroup) {
	user := router.Group("/user")
	{
		user.GET("/list", middlewares.JwtMiddleware(), controllers.GetUsers)
		user.GET("/:id", middlewares.JwtMiddleware(), controllers.GetUser)
		user.PUT("/:id", middlewares.JwtMiddleware(), controllers.Update)
		user.DELETE("/:id", middlewares.JwtMiddleware(), controllers.Delete)
	}
}
