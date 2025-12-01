package controllers

import (
	"health/models"
	db "health/models/db"
	"health/services"
	"health/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Register is an endpoint that creates a new user in the MongoDB database.
// The request body must contain a name, email address and password.
// The email address must be unique. The password is hashed using bcrypt.
// The user is created with the role "user".
// If the user cannot be created, an error is returned.
// Otherwise, a JSON response with the user is sent.
func Register(ctx *gin.Context) {
	var requestBody models.RegisterRequest
	_ = ctx.ShouldBindBodyWith(&requestBody, binding.JSON)

	err := services.CheckUserMail(requestBody.Email)

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	requestBody.Name = strings.TrimSpace(requestBody.Name)
	user, err := services.CreateUser(requestBody.Name, requestBody.Email, requestBody.Password)

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, gin.H{
		"user": user,
	})
}

// Login is an endpoint that verifies the given email and password.
// If the verification succeeds, it generates new access tokens for the user.
// The tokens are then sent in the response as JSON data.
// If the verification fails, it sends a 400 error response with the error message.
func Login(c *gin.Context) {
	var requestBody models.LoginRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	// get user by email
	user, err := services.FindUserByEmail(requestBody.Email)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// check hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// generate new access tokens
	accessToken, refreshToken, err := services.GenerateAccessTokens(user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"user":    user,
		"access":  accessToken.GetResponseJson(),
		"refresh": refreshToken.GetResponseJson(),
	})
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

	// check token validity
	token, err := services.VerifyToken(requestBody.Token, db.TokenTypeRefresh)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := services.FindUserById(token.User)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// delete old token
	err = services.DeleteTokenById(token.ID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, refreshToken, _ := services.GenerateAccessTokens(user)

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"user":    user,
		"access":  accessToken.GetResponseJson(),
		"refresh": refreshToken.GetResponseJson(),
	})
}

// GetAuthProfile is a gin handler that retrieves the user profile of the currently authenticated user.
// The handler expects the user ID to be set in the gin context.
// If the user ID is not set, the handler will send a 400 error response with the error message "cannot get user".
// Otherwise, it will retrieve the user from the database and send a 200 response with the user in the response body.
func GetAuthProfile(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		utils.ErrorResponse(c, http.StatusBadRequest, "cannot get user")
		return
	}
	user, _ := services.FindUserById(userId.(primitive.ObjectID))
	utils.SuccessResponse(c, http.StatusOK, user)
}
