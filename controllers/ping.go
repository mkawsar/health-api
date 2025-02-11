package controllers

import (
	"health/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(ctx *gin.Context) {
	utils.SuccessResponse(ctx, http.StatusOK, "Pong")
}
