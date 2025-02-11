package controllers

import (
	"health/services"
	"health/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUsers is a gin handler that returns a list of all users in the
// database. It paginates the results using the "page" query parameter.
// If the page parameter is not given, it defaults to 0. The response is
// a JSON object with a "users" key that contains the list of users.
// The response also contains "prev" and "next" keys that indicate
// whether there are previous and next pages of results, respectively.
func GetUsers(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	nameFilter := ctx.Query("name")
	users, total, _ := services.GetUSers(ctx.Request.Context(), page, limit, nameFilter)

	utils.PaginatedSuccessResponse(ctx, users, page, limit, total)
}
