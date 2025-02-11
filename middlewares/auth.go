package middlewares

import (
	db "health/models/db"
	"health/services"
	"health/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JwtMiddleware is a middleware function for Gin that checks for a Bearer token
// in the request header. It verifies the token using the VerifyToken service
// and, if valid, sets the user ID in the context. If the token is invalid, it
// sends an unauthorized error response and aborts the request.
func JwtMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "token is required")
			return
		}
		tokenModel, err := services.VerifyToken(token, db.TokenTypeAccess)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
			return
		}

		ctx.Set("userIdHex", tokenModel.User.Hex())
		ctx.Set("userId", tokenModel.User)
		ctx.Next()
	}
}
