package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int         `json:"-"`
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
	StatusCode  int         `json:"-"`
	Success     bool        `json:"success"`
	Message     string      `json:"message"`
	Data        interface{} `json:"data,omitempty"`
	CurrentPage int         `json:"current_page"`
	PerPage     int         `json:"per_page"`
	Total       int64       `json:"total"`
	TotalPages  int64       `json:"total_pages"`
}

// SendResponse sends a JSON response with the given status code and data,
// and marks the response as successful or not.
func (res *Response) SendResponse(c *gin.Context) {
	c.AbortWithStatusJSON(res.StatusCode, res)
}

// SendPaginatedResponse sends a JSON response with the given status code and
// paginated data, and marks the response as successful or not.
func (res *PaginatedResponse) SendPaginatedResponse(c *gin.Context) {
	c.AbortWithStatusJSON(res.StatusCode, res)
}

// SuccessResponse sends a JSON response with the given status code and data,
// and marks the response as successful.
func SuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	response := &Response{
		StatusCode: statusCode,
		Success:    true,
		Message:    "Operation successful",
		Data:       data,
	}
	response.SendResponse(c)
}

// ErrorResponse sends a JSON response with the given status code and message,
// and marks the response as a failure.
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	response := &Response{
		StatusCode: statusCode,
		Success:    false,
		Message:    message,
	}
	response.SendResponse(c)
}

// PaginatedSuccessResponse sends a JSON response with the given data, page,
// per-page limit, and total count, and marks the response as successful.
// The response will also include the total number of pages.
func PaginatedSuccessResponse(c *gin.Context, data interface{}, page int, perPage int, total int64) {
	response := &PaginatedResponse{
		StatusCode:  http.StatusOK,
		Success:     true,
		Message:     "Operation successful",
		Data:        data,
		CurrentPage: page,
		PerPage:     perPage,
		TotalPages:  (total + int64(perPage) - 1) / int64(perPage),
	}
	response.SendPaginatedResponse(c)
}
