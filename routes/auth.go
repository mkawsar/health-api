package routes

import (
	"health/controllers"
	"health/middlewares"
	"health/middlewares/validators"

	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", validators.RegisterValidator(), controllers.Register)
		auth.POST("/login", validators.LoginValidator(), controllers.Login)
		auth.POST("/refresh", validators.RefreshValidator(), controllers.Refresh)
		auth.GET("/profile", middlewares.JwtMiddleware(), middlewares.RoleMiddleware("admin", "user"), controllers.GetAuthProfile)
	}
}
