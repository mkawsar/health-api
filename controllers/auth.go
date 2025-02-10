package controllers

import (
	"health/models"
	"health/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/crypto/bcrypt"
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

func Login(ctx *gin.Context) {
	var requestBody models.LoginRequest
	_ = ctx.ShouldBindBodyWith(&requestBody, binding.JSON)
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// get user by email
	user, err := services.FindUserByEmail(requestBody.Email)

	if err != nil {
		response.Message = err.Error()
		response.SendResponse(ctx)
		return
	}

	// check hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
	if err != nil {
		response.Message = "email and password don't match"
		response.SendResponse(ctx)
		return
	}

	// generate new access tokens
	accessToken, refreshToken, err := services.GenerateAccessTokens(user)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(ctx)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{
		"user": user,
		"token": gin.H{
			"access":  accessToken.GetResponseJson(),
			"refresh": refreshToken.GetResponseJson()},
	}
	response.SendResponse(ctx)
}
