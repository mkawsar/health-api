package validators

import (
	"health/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// RegisterValidator is a middleware that validates the JSON body of a request
// against the models.RegisterRequest struct. If the validation fails, it sends
// a 400 error response with the error message and aborts the request. If the
// validation succeeds, it calls the next handler in the chain.
func RegisterValidator() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var registerRequest models.RegisterRequest
		_ = ctx.ShouldBindBodyWith(&registerRequest, binding.JSON)

		if err := registerRequest.Validate(); err != nil {
			models.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
		ctx.Next()
	}
}

// LoginValidator is a middleware that validates the JSON body of a request
// against the models.LoginRequest struct. If the validation fails, it sends
// a 400 error response with the error message and aborts the request. If the
// validation succeeds, it calls the next handler in the chain.
func LoginValidator() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loginRequest models.LoginRequest
		_ = ctx.ShouldBindBodyWith(&loginRequest, binding.JSON)
		if err := loginRequest.Validate(); err != nil {
			models.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
		ctx.Next()
	}
}

// RefreshValidator is a middleware that validates the JSON body of a request
// against the models.RefreshRequest struct. If the validation fails, it sends
// a 400 error response with the error message and aborts the request. If the
// validation succeeds, it calls the next handler in the chain.
func RefreshValidator() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var refreshRequest models.RefreshRequest
		_ = ctx.ShouldBindBodyWith(&refreshRequest, binding.JSON)
		if err := refreshRequest.Validate(); err != nil {
			models.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
		ctx.Next()
	}
}
