package controllers

import (
	"health/services"
	"health/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary      Get a list of users
// @Description  Get a paginated list of users
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.PaginatedResponse
// @Param        page      query     int     false  "Page number"     default(1)
// @Param        limit  query     int     false  "Items per page"  default(10)
// @Param        name  query     string     false  "sorted by name"
// @Router       /v1/user/list [get]
// @Security     ApiKeyAuth
func GetUsers(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	nameFilter := ctx.Query("name")
	users, total, _ := services.GetUSers(ctx.Request.Context(), page, limit, nameFilter)

	utils.PaginatedSuccessResponse(ctx, users, page, limit, total)
}

func GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	user, _ := services.GetUser(objectID)
	utils.SuccessResponse(ctx, http.StatusOK, user)
}
