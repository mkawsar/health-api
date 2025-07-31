package controllers

import (
	"health/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDoctors(ctx *gin.Context) {
	utils.SuccessResponse(ctx, http.StatusOK, nil)
}
