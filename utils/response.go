package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PaginatedResponse struct {
	Success     bool        `json:"success"`
	Message     string      `json:"message"`
	Data        interface{} `json:"data"`
	CurrentPage int         `json:"current_page"`
	PerPage     int         `json:"per_page"`
	Total       int64       `json:"total"`
	TotalPages  int64       `json:"total_pages"`
}

// SuccessResponse sends a JSON response with the given status code and data,
// and marks the response as successful.
func SuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: "Operation successful",
		Data:    data,
	})
}

// ErrorResponse sends a JSON response with the given status code and message,
// and marks the response as a failure.
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Data:    nil,
	})
}

// PaginatedSuccessResponse sends a JSON response with the given data, page,
// per-page limit, and total count, and marks the response as successful.
// The response will also include the total number of pages.
func PaginatedSuccessResponse(c *gin.Context, data interface{}, page int, perPage int, total int64) {
	totalPages := (total + int64(perPage) - 1) / int64(perPage)
	c.JSON(http.StatusOK, PaginatedResponse{
		Success:     true,
		Message:     "Operation successful",
		Data:        data,
		CurrentPage: page,
		PerPage:     perPage,
		Total:       total,
		TotalPages:  totalPages,
	})
}
