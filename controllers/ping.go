package controllers

import (
	"health/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(ctx *gin.Context) {
	response := &models.Response {
		StatusCode: http.StatusOK,
		Success: true,
		Message: "pong",
	}

	response.SendResponse(ctx)
}
