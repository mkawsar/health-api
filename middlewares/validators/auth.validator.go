package validators

import (
	"health/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

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
