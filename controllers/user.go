package controllers

import (
	"health/services"
	"health/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetUsers is a gin handler that returns a list of all users in the
// database. It paginates the results using the "page" query parameter.
// If the page parameter is not given, it defaults to 0. The response is
// a JSON object with a "users" key that contains the list of users.
// The response also contains "prev" and "next" keys that indicate
// whether there are previous and next pages of results, respectively.

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
