package routes

import (
	"health/controllers"
	"health/middlewares"

	"github.com/gin-gonic/gin"
)

func DoctorRoute(router *gin.RouterGroup) {
	doctor := router.Group("/doctor")
	{
		doctor.GET("/list", middlewares.JwtMiddleware(), controllers.GetDoctors)
	}
}