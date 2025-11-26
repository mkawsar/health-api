package validators

import (
	"health/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// PathIdValidator is a middleware that validates the "id" path parameter
// as a numeric ID. If the validation fails, it sends a 400 error response
// with the error message and aborts the request. If the validation succeeds,
// it calls the next handler in the chain.
func PathIdValidator() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		err := validation.Validate(id, validation.Required, is.Int)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id: "+id)
			return
		}

		ctx.Next()
	}
}
