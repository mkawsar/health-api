package middlewares

import (
	db "health/models/db"
	"health/services"
	"health/utils"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)


// JwtMiddleware is a middleware that verifies a JWT token from the Authorization header
// and sets the userId, userIdHex, and role fields in the gin context.
// If the token is invalid or the user associated with the token cannot be found,
// it sends an unauthorized error response and aborts the request.
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

		user, _ := services.FindUserById(tokenModel.User)
		if user == nil {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "user not found")
			return
		}

		ctx.Set("userIdHex", tokenModel.User.Hex())
		ctx.Set("userId", tokenModel.User)
		ctx.Set("role", user.Role)
		ctx.Next()
	}
}

// RoleMiddleware is a middleware function for Gin that checks the user's role
// against a given set of allowed roles. If the user's role is not in the allowed
// roles, it sends a forbidden error response and aborts the request. Otherwise, it
// calls the next handler in the chain.
func RoleMiddleware(allowedRole ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRoleIfc, exists := ctx.Get("role")
		if !exists {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "No role found")
			return
		}

		userRole := userRoleIfc.(string)
		if sort.SearchStrings(allowedRole, userRole) < len(allowedRole) {
			ctx.Next()
			return
		}
		// If the user role is not in the allowed roles, return a forbidden error
		// and abort the request
		utils.ErrorResponse(ctx, http.StatusForbidden, "You don't have permission to access this resource")
		ctx.Abort()
	}
}
