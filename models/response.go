package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int            `json:"-"`
	Success    bool           `json:"success"`
	Message    string         `json:"message,omitempty"`
	Data       map[string]any `json:"data,omitempty"`
}

// SendResponse is a helper function to send a gin response and abort the context.
// It will use the StatusCode of the Response struct to set the HTTP status code.
func (res *Response) SendResponse(c *gin.Context) {
	c.AbortWithStatusJSON(res.StatusCode, res)
}

// SendResponseData sends a gin response and aborts the context.
// It will use the given gin.H as the response data.
// The HTTP status code will be set to 200.
func SendResponseData(c *gin.Context, data gin.H) {
	response := &Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Data:       data,
	}
	response.SendResponse(c)
}

func SendErrorResponse(c *gin.Context, status int, msg string) {
	response := &Response{
		StatusCode: status,
		Success:    false,
		Message:    msg,
	}
	response.SendResponse(c)
}
