package middlewares

import (
	"health/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AppRecovery returns a gin error handler that will catch panics and write a 500 response to the client.
// If the panic value is a string, it will be used as the error message in the response.
// Otherwise, the response will only contain a generic error message.
func AppRecovery() func(ctx *gin.Context, recovered interface{}) {
	return func(ctx *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			models.SendErrorResponse(ctx, http.StatusInternalServerError, err)
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"success": false})
	}
}
