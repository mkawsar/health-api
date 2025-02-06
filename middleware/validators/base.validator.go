package validators

import (
	"health/models"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func PathIdValidator() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		err := validation.Validate(id, is.MongoID)
		if err != nil {
			models.SendErrorResponse(ctx, http.StatusBadRequest, "invalid id: "+id)
			return
		}

		ctx.Next()
	}
}
