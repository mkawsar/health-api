package controllers

import (
	"health/services"
	"health/utils"
	"health/utils/requests"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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

// @Summary      Get a user by ID
// @Description  Get a user by the given ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Router       /v1/user/{id} [get]
// @Security     ApiKeyAuth
func GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid user id")
		return
	}
	user, _ := services.GetUser(uint(userId))
	utils.SuccessResponse(ctx, http.StatusOK, user)
}

// @Summary      Update a user
// @Description  Update a user by the given ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.Response
// @Param        id   path      string     true  "User ID"
// @Param        UserRequest  body      requests.UserRequest  true  "User details"
// @Router       /v1/user/{id} [patch]
// @Security     ApiKeyAuth
func Update(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid user id")
		return
	}
	var request requests.UserRequest
	_ = ctx.ShouldBindBodyWith(&request, binding.JSON)

	err = services.UpdateUser(uint(userId), &request)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "User updated successfully")
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid user id")
		return
	}
	
	err = services.DeleteUser(uint(userId))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	
	utils.SuccessResponse(ctx, http.StatusOK, "User deleted successfully")
}
