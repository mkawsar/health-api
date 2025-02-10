package controllers

import (
	"health/models"
	db "health/models/db"
	"health/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// Login is an endpoint that verifies the given email and password.
// If the verification succeeds, it generates new access tokens for the user.
// The tokens are then sent in the response as JSON data.
// If the verification fails, it sends a 400 error response with the error message.
func Login(c *gin.Context) {
	var requestBody models.LoginRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// get user by email
	user, err := services.FindUserByEmail(requestBody.Email)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// check hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
	if err != nil {
		response.Message = "email and password don't match"
		response.SendResponse(c)
		return
	}

	// generate new access tokens
	accessToken, refreshToken, err := services.GenerateAccessTokens(user)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
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
	response.SendResponse(c)
}

// Refresh is a gin handler that refreshes an access token using a refresh token.
// The handler expects a JSON body with a "token" field that contains the refresh token.
// The handler will verify the token, find the associated user, delete the old token
// and generate new access tokens. If the token is invalid, the associated user cannot
// be found, the old token cannot be deleted, or the new tokens cannot be generated,
// the handler will send a 400 error response with the error message. Otherwise, it will
// send a 200 response with the user and the new tokens in the response body.
func Refresh(c *gin.Context) {
	var requestBody models.RefreshRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// check token validity
	token, err := services.VerifyToken(requestBody.Token, db.TokenTypeRefresh)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	user, err := services.FindUserById(token.User)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// delete old token
	err = services.DeleteTokenById(token.ID)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	accessToken, refreshToken, _ := services.GenerateAccessTokens(user)
	
	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{
		"user": user,
		"token": gin.H{
			"access":  accessToken.GetResponseJson(),
			"refresh": refreshToken.GetResponseJson()},
	}
	response.SendResponse(c)
}

func GetAuthProfile(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}
	userId, exists := c.Get("userId")
	
	if !exists {
		response.Message = "cannot get user"
		response.SendResponse(c)
		return
	}
	user, _ := services.FindUserById(userId.(primitive.ObjectID))
	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{
		"user": user,
	}
	response.SendResponse(c)
}
