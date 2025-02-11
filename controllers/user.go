package controllers

import (
	"health/models"
	"health/services"
	"net/http"
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
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	nameFilter := ctx.Query("name")
	users, totalUsers, _ := services.GetUSers(ctx.Request.Context(), page, limit, nameFilter)

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{
		"users": users,
		"pagination": gin.H{
			"current_page": page,
			"per_page":     limit,
			"total_users":  totalUsers,
			"total_pages":  (totalUsers + int64(limit) - 1) / int64(limit),
		},
	}
	response.SendResponse(ctx)
}
