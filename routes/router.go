package routes

import (
	"health/middlewares"
	"health/models"
	"health/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	r := gin.New()
	initRoutes(r)

	r.Use(gin.LoggerWithWriter(middlewares.LogWriter()))
	r.Use(gin.CustomRecovery(middlewares.AppRecovery()))
	r.Use(middlewares.CORSMiddleware())
	v1 := r.Group("/v1")
	{
		PingRoute(v1)
	}
	return r
}

func initRoutes(r *gin.Engine) {
	_ = r.SetTrustedProxies(nil)
	r.RedirectTrailingSlash = false
	r.HandleMethodNotAllowed = true
	r.NoRoute(func(ctx *gin.Context) {
		models.SendErrorResponse(ctx, http.StatusNotFound, ctx.Request.RequestURI+" not found")
	})
	r.NoMethod(func(ctx *gin.Context) {
		models.SendErrorResponse(ctx, http.StatusMethodNotAllowed, ctx.Request.Method+" is not allowed here")
	})
}

func InitGin() {
	gin.DisableConsoleColor()
	gin.SetMode(services.Config.Mode)
}
