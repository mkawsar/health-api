package routes

import (
	"health/middlewares"
	"health/services"
	"health/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// New returns a new gin.Engine instance with routes and middlewares set up.
func New() *gin.Engine {
	r := gin.New()
	initRoutes(r)

	r.Use(gin.LoggerWithWriter(middlewares.LogWriter()))
	r.Use(gin.CustomRecovery(middlewares.AppRecovery()))
	r.Use(middlewares.CORSMiddleware())
	v1 := r.Group("/v1")
	{
		PingRoute(v1)
		AuthRoute(v1)
		UserRoute(v1)
	}
	return r
}

// initRoutes sets up the router to redirect trailing slashes, handle
// method-not-allowed and not-found requests, and sets up custom 404 handlers.
func initRoutes(r *gin.Engine) {
	_ = r.SetTrustedProxies(nil)
	r.RedirectTrailingSlash = false
	r.HandleMethodNotAllowed = true
	r.NoRoute(func(ctx *gin.Context) {
		utils.ErrorResponse(ctx, http.StatusNotFound, ctx.Request.RequestURI+" not found")
	})
	r.NoMethod(func(ctx *gin.Context) {
		utils.ErrorResponse(ctx, http.StatusNotFound, ctx.Request.Method+" is not allowed here")
	})
}

// InitGin configures the Gin mode to the mode specified in the configuration.
// It is called once during application startup.
func InitGin() {
	gin.DisableConsoleColor()
	gin.SetMode(services.Config.Mode)
}
