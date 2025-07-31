package controllers

import (
	"health/services"
	"health/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetDoctors(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	nameFilter := ctx.Query("name")
	users, total, _ := services.GetDoctors(ctx.Request.Context(), page, limit, nameFilter)

	utils.PaginatedSuccessResponse(ctx, users, page, limit, total)
}
