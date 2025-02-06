package controllers

import (
	"health/models"
	"health/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Register(ctx *gin.Context) {
	var requestBody models.RegisterRequest
	_ = ctx.ShouldBindBodyWith(&requestBody, binding.JSON)
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	err := services.CheckUserMail(requestBody.Email)

	if err != nil {
		response.Message = err.Error()
		response.SendResponse(ctx)
		return
	}

	requestBody.Name = strings.TrimSpace(requestBody.Name)
	user, err := services.CreateUser(requestBody.Name, requestBody.Email, requestBody.Password)

	if err != nil {
		response.Message = err.Error()
		response.SendResponse(ctx)
		return
	}

	response.Data = gin.H{
		"user": user,
	}

	response.SendResponse(ctx)
}
